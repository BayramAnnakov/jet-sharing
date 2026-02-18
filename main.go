package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// module github.com/jet-sharing/scooter-api

// Scooter represents an electric scooter in the fleet.
type Scooter struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	BatteryLevel   int     `json:"battery_level"` // percentage 0-100
	Status         string  `json:"status"`        // "available", "in_use", "maintenance"
	PricePerMinute float64 `json:"price_per_minute"`
}

var (
	mu       sync.RWMutex
	scooters = map[string]*Scooter{
		"sc-1001": {ID: "sc-1001", Name: "Av. Paulista & Rua Augusta", Latitude: -23.5613, Longitude: -46.6560, BatteryLevel: 87, Status: "available", PricePerMinute: 0.50},
		"sc-1002": {ID: "sc-1002", Name: "Praça da Sé", Latitude: -23.5503, Longitude: -46.6340, BatteryLevel: 42, Status: "available", PricePerMinute: 0.50},
		"sc-1003": {ID: "sc-1003", Name: "Parque Ibirapuera", Latitude: -23.5874, Longitude: -46.6576, BatteryLevel: 15, Status: "maintenance", PricePerMinute: 0.50},
		"sc-1004": {ID: "sc-1004", Name: "Pinheiros", Latitude: -23.5670, Longitude: -46.6914, BatteryLevel: 95, Status: "in_use", PricePerMinute: 0.50},
	}
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	r.Get("/health", handleHealthCheck)
	r.Route("/api/scooters", func(r chi.Router) {
		r.Get("/", handleListScooters)
		r.Get("/{id}", handleGetScooter)
		r.Post("/{id}/unlock", handleUnlockScooter)
		r.Post("/{id}/lock", handleLockScooter)
	})

	addr := ":8080"
	slog.Info("starting scooter API server", "addr", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleListScooters(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]*Scooter, 0, len(scooters))
	for _, s := range scooters {
		result = append(result, s)
	}

	slog.Info("listing scooters", "count", len(result))
	writeJSON(w, http.StatusOK, result)
}

func handleGetScooter(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	scooter, err := findScooter(r.Context(), id)
	if err != nil {
		slog.Warn("scooter not found", "id", id)
		writeError(w, http.StatusNotFound, "scooter not found")
		return
	}

	writeJSON(w, http.StatusOK, scooter)
}

func handleUnlockScooter(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	mu.Lock()
	defer mu.Unlock()

	scooter, ok := scooters[id]
	if !ok {
		writeError(w, http.StatusNotFound, "scooter not found")
		return
	}

	if scooter.Status != "available" {
		slog.Warn("unlock rejected", "id", id, "status", scooter.Status)
		writeError(w, http.StatusConflict, fmt.Sprintf("scooter is %s, cannot unlock", scooter.Status))
		return
	}

	if scooter.BatteryLevel < 10 {
		writeError(w, http.StatusBadRequest, "battery too low for ride")
		return
	}

	// TODO: verify payment method before unlocking
	// paymentOk, err := verifyPayment(r.Context(), userID)

	scooter.Status = "in_use"
	slog.Info("scooter unlocked", "id", id, "battery", scooter.BatteryLevel)
	writeJSON(w, http.StatusOK, scooter)
}

func handleLockScooter(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	mu.Lock()
	defer mu.Unlock()

	scooter, ok := scooters[id]
	if !ok {
		writeError(w, http.StatusNotFound, "scooter not found")
		return
	}

	if scooter.Status != "in_use" {
		writeError(w, http.StatusConflict, fmt.Sprintf("scooter is %s, cannot lock", scooter.Status))
		return
	}

	scooter.Status = "available"
	slog.Info("scooter locked", "id", id)
	writeJSON(w, http.StatusOK, scooter)
}

// findScooter retrieves a scooter by ID from the in-memory store.
func findScooter(_ context.Context, id string) (*Scooter, error) {
	mu.RLock()
	defer mu.RUnlock()

	s, ok := scooters[id]
	if !ok {
		return nil, fmt.Errorf("findScooter: %w", ErrScooterNotFound)
	}
	return s, nil
}

var ErrScooterNotFound = fmt.Errorf("scooter not found")

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// writeJSON sends a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

// writeError sends a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
