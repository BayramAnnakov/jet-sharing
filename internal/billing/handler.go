package billing

import (
	"fmt"
	"log/slog"
	"math"
)

// Config holds per-zone billing configuration.
type Config struct {
	Zone            string  `json:"zone"`
	BaseFare        float64 `json:"base_fare"`
	PerMinuteRate   float64 `json:"per_minute_rate"`
	PerKmRate       float64 `json:"per_km_rate"`
	MinFare         float64 `json:"min_fare"`
	SurgeMultiplier float64 `json:"surge_multiplier"`
	Currency        string  `json:"currency"`
}

// FareResult contains the calculated fare breakdown.
type FareResult struct {
	BaseFare    float64 `json:"base_fare"`
	TimeFare    float64 `json:"time_fare"`
	DistFare    float64 `json:"distance_fare"`
	SurgeAmount float64 `json:"surge_amount"`
	TotalFare   float64 `json:"total_fare"`
	Currency    string  `json:"currency"`
}

// defaultConfigs holds billing configuration per zone.
var defaultConfigs = map[string]*Config{
	"baku-central": {Zone: "baku-central", BaseFare: 2.50, PerMinuteRate: 0.50, PerKmRate: 1.20, MinFare: 5.00, SurgeMultiplier: 1.0, Currency: "AZN"},
	"baku-south":   {Zone: "baku-south", BaseFare: 2.00, PerMinuteRate: 0.45, PerKmRate: 1.00, MinFare: 4.50, SurgeMultiplier: 1.0, Currency: "AZN"},
}

// LoadBillingConfig retrieves the billing configuration for a given zone.
func LoadBillingConfig(zone string) (*Config, error) {
	cfg, ok := defaultConfigs[zone]
	if !ok {
		return nil, fmt.Errorf("billing: unknown zone %q", zone)
	}

	slog.Info("BillingConfigLoaded",
		"zone", zone,
		"base_fare", cfg.BaseFare,
		"per_minute", cfg.PerMinuteRate,
		"per_km", cfg.PerKmRate,
	)

	return cfg, nil
}

// CalculateFare computes the ride fare based on distance, duration, and zone config.
func CalculateFare(distanceKm float64, durationMin float64, cfg *Config) *FareResult {
	timeFare := durationMin * cfg.PerMinuteRate
	distFare := distanceKm * cfg.PerKmRate
	subtotal := cfg.BaseFare + timeFare + distFare

	surgeAmount := 0.0
	if cfg.SurgeMultiplier > 1.0 {
		surgeAmount = subtotal * (cfg.SurgeMultiplier - 1.0)
	}

	total := subtotal + surgeAmount
	if total < cfg.MinFare {
		total = cfg.MinFare
	}

	// Round to 2 decimal places
	total = math.Round(total*100) / 100

	result := &FareResult{
		BaseFare:    cfg.BaseFare,
		TimeFare:    timeFare,
		DistFare:    distFare,
		SurgeAmount: surgeAmount,
		TotalFare:   total,
		Currency:    cfg.Currency,
	}

	slog.Info("FareCalculated",
		"distance_km", distanceKm,
		"duration_min", durationMin,
		"total_fare", result.TotalFare,
		"currency", result.Currency,
	)

	return result
}

// PerMinuteStep logs a per-minute billing tick during an active ride.
func PerMinuteStep(orderID string, minuteIndex int, runningTotal float64, cfg *Config) float64 {
	step := cfg.PerMinuteRate * cfg.SurgeMultiplier
	newTotal := runningTotal + step

	slog.Info("PerMinuteStep",
		"order_id", orderID,
		"minute", minuteIndex,
		"step_amount", step,
		"running_total", math.Round(newTotal*100)/100,
	)

	return newTotal
}
