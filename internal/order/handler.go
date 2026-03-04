package order

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// Store is the interface for order persistence.
type Store interface {
	GetOrder(ctx context.Context, orderID string) (*Order, error)
	CreateOrder(ctx context.Context, order *Order) error
	UpdateOrder(ctx context.Context, order *Order) error
}

// Order represents a ride order.
type Order struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ScooterID   string    `json:"scooter_id"`
	StatusID    int       `json:"status_id"`
	FareAmount  float64   `json:"fare_amount"`
	Currency    string    `json:"currency"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

// Handler holds dependencies for order HTTP handlers.
type Handler struct {
	store Store
	mu    sync.RWMutex
}

// NewHandler creates a new order handler.
func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

// HandleCreateOrder creates a new ride order.
func (h *Handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    string `json:"user_id"`
		ScooterID string `json:"scooter_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	order := &Order{
		ID:        generateOrderID(),
		UserID:    req.UserID,
		ScooterID: req.ScooterID,
		StatusID:  StatusCreated,
		Currency:  "AZN",
		CreatedAt: time.Now(),
	}

	if err := h.store.CreateOrder(r.Context(), order); err != nil {
		slog.Error("failed to create order", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	slog.Info("OrderCreated",
		"order_id", order.ID,
		"user_id", order.UserID,
		"scooter_id", order.ScooterID,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// HandleGetOrderStatus returns the current status for an order.
// Called by the mobile app to poll for status updates.
func (h *Handler) HandleGetOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	pollSeq := r.URL.Query().Get("poll_seq")

	// BUG: no rate limiting, no caching — every poll hits the store and logs.
	// The mobile app polls every 2s, causing 143 polls in 4.6 minutes for a single order.
	order, err := h.store.GetOrder(r.Context(), orderID)
	if err != nil {
		slog.Warn("order not found", "order_id", orderID)
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	slog.Info("GetOrderStatus",
		"order_id", orderID,
		"status_id", order.StatusID,
		"poll_seq", pollSeq,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// HandleOrderCompleted is called when an order reaches terminal state.
func (h *Handler) HandleOrderCompleted(ctx context.Context, order *Order) {
	order.CompletedAt = time.Now()
	if err := h.store.UpdateOrder(ctx, order); err != nil {
		slog.Error("failed to update completed order", "order_id", order.ID, "error", err)
		return
	}

	slog.Info("OrderCompleted",
		"order_id", order.ID,
		"fare_amount", order.FareAmount,
		"duration_min", time.Since(order.CreatedAt).Minutes(),
	)
}

func generateOrderID() string {
	return "ORD-" + time.Now().Format("20060102-150405")
}
