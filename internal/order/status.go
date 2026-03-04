package order

// Order status codes — matches orders.status in the database.
// See DATABASE.md for the full reference.
const (
	StatusPending          = 1
	StatusCreated          = 2
	StatusScooterAssigned  = 3
	StatusRouteCalculated  = 4
	StatusRideActive       = 5
	StatusRidePaused       = 6
	StatusRideEnding       = 7
	StatusPaymentPending   = 10
	StatusPaymentProcessing = 11
	StatusPaymentCompleted = 12
	StatusPaymentFailed    = 13
	StatusCompleted        = 14
	StatusCancelled        = 15
	StatusRefunded         = 16
)
// Codes 8-9 are reserved for future ride states.

// allowedTransitions defines which status transitions are valid.
// Key is the current status, value is a slice of allowed next statuses.
var allowedTransitions = map[int][]int{
	StatusPending:          {StatusCreated, StatusCancelled},
	StatusCreated:          {StatusScooterAssigned, StatusCancelled},
	StatusScooterAssigned:  {StatusRouteCalculated, StatusCancelled},
	StatusRouteCalculated:  {StatusRideActive, StatusCancelled},
	StatusRideActive:       {StatusRidePaused, StatusRideEnding, StatusCompleted}, // fast-path for pre-paid orders
	StatusRidePaused:       {StatusRideActive, StatusRideEnding},
	StatusRideEnding:       {StatusPaymentPending},
	StatusPaymentPending:   {StatusPaymentProcessing, StatusPaymentFailed},
	StatusPaymentProcessing: {StatusPaymentCompleted, StatusPaymentFailed},
	StatusPaymentCompleted: {StatusCompleted},
	StatusPaymentFailed:    {StatusPaymentPending, StatusCancelled},
	StatusCompleted:        {StatusRefunded},
	StatusCancelled:        {StatusRefunded},
}

// IsValidTransition checks whether transitioning from `from` to `to` is allowed.
func IsValidTransition(from, to int) bool {
	allowed, ok := allowedTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}
