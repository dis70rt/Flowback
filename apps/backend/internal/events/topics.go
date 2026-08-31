package events

// TopicWebhookReceived is the single unified Asynq queue topic.
// All Razorpay webhook events are pushed to one queue and routed by event type in the worker.
const TopicWebhookReceived = "webhook:received"

// Broadcast channel names used by the broadcaster to publish decline decisions.
// Downstream services (retry engine, notification service, AI agent, etc.) subscribe to these.
const (
	TopicDeclineHard    = "decline:hard"    // Payment is unrecoverable — customer must provide new method
	TopicDeclineSoft    = "decline:soft"    // Payment is retryable — schedule a retry attempt
	TopicDeclineUnknown = "decline:unknown" // Cannot be classified — forward to AI agent for decision
)

// Razorpay Webhook Event Types.
// Verified against the official Razorpay docs (April 2026):
// - https://razorpay.com/docs/build/llm-docs/webhooks.md
// - https://razorpay.com/docs/build/llm-docs/payments/subscriptions/subscribe-to-webhooks.md

const (
	// --- Payment Events ---

	// EventPaymentAuthorized fires when a payment is authorized (funds reserved but not yet captured).
	// Critical for catching late-authorizations via webhooks.
	EventPaymentAuthorized = "payment.authorized"

	// EventPaymentCaptured fires when a payment is successfully captured (funds moving to settlement).
	EventPaymentCaptured = "payment.captured"

	// EventPaymentFailed fires when a payment attempt fails.
	// This is the primary event for decline classification.
	EventPaymentFailed = "payment.failed"

	// --- Order Events ---

	// EventOrderPaid fires when an order is fully paid.
	// Payload contains both the order and payment entity, making it richer than payment.captured alone.
	EventOrderPaid = "order.paid"

	// --- Refund Events ---

	EventRefundCreated      = "refund.created"
	EventRefundProcessed    = "refund.processed"
	EventRefundFailed       = "refund.failed"
	EventRefundSpeedChanged = "refund.speed_changed"

	// --- Subscription Events ---
	// Verified against: https://razorpay.com/docs/build/llm-docs/payments/subscriptions/subscribe-to-webhooks.md
	//
	// subscription.authenticated -> first payment made (auth/upfront/plan amount)
	// subscription.pending       -> a charge failed; Razorpay is retrying
	// subscription.halted        -> all retries exhausted; customer must act
	// subscription.cancelled     -> subscription cancelled and moved to cancelled state

	EventSubscriptionAuthenticated = "subscription.authenticated"
	EventSubscriptionActivated     = "subscription.activated"
	EventSubscriptionCharged       = "subscription.charged"
	EventSubscriptionCompleted     = "subscription.completed"
	EventSubscriptionUpdated       = "subscription.updated"
	EventSubscriptionPending       = "subscription.pending"
	EventSubscriptionHalted        = "subscription.halted"
	EventSubscriptionCancelled     = "subscription.cancelled"
	EventSubscriptionPaused        = "subscription.paused"
	EventSubscriptionResumed       = "subscription.resumed"

	// --- Settlement Events ---

	EventSettlementProcessed = "settlement.processed"
)

// WebhookPayload is the single, unified firehose payload pushed to the Asynq queue.
// The API handler does zero parsing — it just peeks at the event string and stores raw bytes.
// All parsing is done inside the worker so the API remains a dumb ingestion layer.
type WebhookPayload struct {
	Event   string `json:"event"`
	RawJSON []byte `json:"raw_json"`
}
