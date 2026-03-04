package order

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Lifecycle manages order state transitions.
type Lifecycle struct {
	store Store
}

// NewLifecycle creates a new lifecycle manager.
func NewLifecycle(store Store) *Lifecycle {
	return &Lifecycle{store: store}
}

// TransitionStatus validates and applies a status transition for an order.
func (l *Lifecycle) TransitionStatus(ctx context.Context, orderID string, newStatus int) error {
	order, err := l.store.GetOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("lifecycle: get order %s: %w", orderID, err)
	}

	oldStatus := order.StatusID

	if !IsValidTransition(oldStatus, newStatus) {
		slog.Warn("invalid status transition",
			"order_id", orderID,
			"from", oldStatus,
			"to", newStatus,
		)
		return fmt.Errorf("lifecycle: invalid transition from %d to %d", oldStatus, newStatus)
	}

	order.StatusID = newStatus
	if err := l.store.UpdateOrder(ctx, order); err != nil {
		return fmt.Errorf("lifecycle: update order %s: %w", orderID, err)
	}

	slog.Info("StatusUpdate",
		"order_id", orderID,
		"from_status", oldStatus,
		"to_status", newStatus,
		"timestamp", time.Now().Unix(),
	)

	return nil
}

// StartRide transitions an order to ride_active state.
// Called after scooter assignment and route calculation are complete.
func (l *Lifecycle) StartRide(ctx context.Context, orderID string) error {
	order, err := l.store.GetOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("lifecycle: get order %s: %w", orderID, err)
	}

	// BUG: only checks order status, does NOT verify the scooter was actually unlocked.
	// If there's a delay between route calculation and scooter unlock (depends on NATS
	// message delivery), the ride can start before the scooter is ready — explaining
	// the 16-minute gap between route_calculated and ride_active in some orders.
	if order.StatusID != StatusRouteCalculated {
		return fmt.Errorf("lifecycle: cannot start ride, order %s is in status %d (expected %d)",
			orderID, order.StatusID, StatusRouteCalculated)
	}

	order.StatusID = StatusRideActive
	if err := l.store.UpdateOrder(ctx, order); err != nil {
		return fmt.Errorf("lifecycle: update order %s: %w", orderID, err)
	}

	slog.Info("RideStarted",
		"order_id", orderID,
		"scooter_id", order.ScooterID,
		"started_at", time.Now().Unix(),
	)

	return nil
}
