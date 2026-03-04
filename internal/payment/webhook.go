package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"jetsharing/internal/order"
)

// StripeEvent represents an incoming Stripe webhook event.
type StripeEvent struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Data struct {
		Object struct {
			PaymentIntentID string  `json:"payment_intent"`
			Amount          float64 `json:"amount"`
			Currency        string  `json:"currency"`
			Status          string  `json:"status"`
			OrderID         string  `json:"metadata_order_id"`
		} `json:"object"`
	} `json:"data"`
}

// WebhookHandler processes Stripe payment webhooks.
type WebhookHandler struct {
	orderStore order.Store
	taskMgr    TaskManager
}

// TaskManager is the interface for managing async tasks.
type TaskManager interface {
	DeleteTask(taskID string) error
	CancelProcessing(taskID string) error
}

// NewWebhookHandler creates a new webhook handler.
func NewWebhookHandler(orderStore order.Store, taskMgr TaskManager) *WebhookHandler {
	return &WebhookHandler{
		orderStore: orderStore,
		taskMgr:    taskMgr,
	}
}

// HandlePaymentWebhook receives and processes Stripe webhook events.
func (h *WebhookHandler) HandlePaymentWebhook(w http.ResponseWriter, r *http.Request) {
	var event StripeEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		slog.Error("failed to decode webhook", "error", err)
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	slog.Info("PaymentCallback",
		"event_id", event.ID,
		"event_type", event.Type,
		"order_id", event.Data.Object.OrderID,
		"amount", event.Data.Object.Amount,
	)

	switch event.Type {
	case "payment_intent.succeeded":
		if err := h.confirmPayment(r.Context(), event); err != nil {
			slog.Error("payment confirmation failed", "error", err, "order_id", event.Data.Object.OrderID)
			http.Error(w, "processing failed", http.StatusInternalServerError)
			return
		}
	case "payment_intent.payment_failed":
		slog.Warn("payment failed",
			"order_id", event.Data.Object.OrderID,
			"status", event.Data.Object.Status,
		)
	default:
		slog.Debug("unhandled webhook event", "type", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

// confirmPayment marks a payment as confirmed and completes the order.
func (h *WebhookHandler) confirmPayment(ctx context.Context, event StripeEvent) error {
	orderID := event.Data.Object.OrderID

	o, err := h.orderStore.GetOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("confirmPayment: get order: %w", err)
	}

	slog.Info("PaymentConfirmed",
		"order_id", orderID,
		"amount", event.Data.Object.Amount,
		"currency", event.Data.Object.Currency,
	)

	// TODO: use lifecycle.TransitionStatus instead of direct update
	// This bypasses the state machine — jumps from ride_active (5) straight to
	// completed (14) without going through ride_ending (7), payment_pending (10),
	// payment_processing (11), payment_completed (12) intermediate states.
	o.StatusID = order.StatusCompleted
	o.FareAmount = event.Data.Object.Amount
	o.CompletedAt = time.Now()

	if err := h.orderStore.UpdateOrder(ctx, o); err != nil {
		return fmt.Errorf("confirmPayment: update order: %w", err)
	}

	// Clean up the payment check task
	taskID := fmt.Sprintf("payment-check-%s", orderID)
	h.taskMgr.DeleteTask(taskID) // BUG: should use CancelProcessing for active tasks

	return nil
}
