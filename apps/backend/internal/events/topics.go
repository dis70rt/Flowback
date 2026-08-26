package events

const (
	TopicPaymentFailed = "payment:failed"
	TopicSubscriptionHalted = "subscription:halted"
	TopicSubscriptionPending = "subscription:pending"
)

type PaymentFailedPayload struct {
	PaymentID     string `json:"payment_id"`
	ErrorCode     string `json:"error_code"`
	ErrorDesc     string `json:"error_desc"`
	CustomerID    string `json:"customer_id,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
}

type SubscriptionHaltedPayload struct {
	SubscriptionID string `json:"subscription_id"`
	CustomerID     string `json:"customer_id,omitempty"`
	PlanID         string `json:"plan_id,omitempty"`
}
