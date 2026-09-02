package decline

import (
	"encoding/json"
)

// DeclineType represents whether a payment failure is recoverable or not.
type DeclineType string

const (
	Soft    DeclineType = "soft"
	Hard    DeclineType = "hard"
	Unknown DeclineType = "unknown" // Sent to AI agent for final decision
)

// ClassificationResult holds the extracted data and the final decline decision.
type ClassificationResult struct {
	Type        DeclineType
	PaymentID   string
	ErrorCode   string
	ErrorReason string
}

// softDeclineCodes are transient failures or failures the customer can fix.
// Retrying automatically or after a short delay has a realistic chance of success.
//
// Sources verified against official Razorpay docs:
//   - https://razorpay.com/docs/build/llm-docs/payments/payments/failure-analysis.md
//   - https://razorpay.com/docs/build/llm-docs/errors/payments/payment-methods-error-parameters.md
var softDeclineCodes = map[string]bool{
	// Customer drop-offs — customer can fix and retry immediately
	"payment_cancelled":               true, // Explicit cancel; show a "try again" prompt
	"payment_collect_request_expired": true, // UPI collect timeout (5-10 min window missed)
	"insufficient_funds":              true, // Docs: "retry with different card or method"
	"insufficient_fund":               true, // Alternate key sometimes seen in webhooks
	"low_balance":                     true, // UPI-specific insufficient funds

	// Transaction limit errors — customer can retry after limit resets or use another method
	"transaction_limit_exceeded":           true,
	"transaction_frequency_limit_exceeded": true,

	// Auth errors that can be fixed by the customer retrying with correct info
	"invalid_otp":   true, // Example shown in official error structure docs
	"incorrect_pin": true,
	"pin_not_set":   true,

	// Transient bank / network / gateway errors — retry automatically after a delay
	"payment_timed_out":       true, // Docs: "retry and complete within time"
	"bank_technical_error":    true, // Docs: "try using another bank account or method"
	"gateway_technical_error": true, // Docs: "retry after some time"
	"server_error":            true, // Razorpay-side transient error
	"payment_failed":          true, // Generic bank/wallet error with no specific code — retry eligible

	// UPI-specific transient failures
	"upi_app_not_available": true, // App temporarily unavailable; customer can retry
}

// hardDeclineCodes are permanent failures. Auto-retry will never succeed.
// The customer MUST provide a new payment method or take action with their bank.
//
// Sources verified against official Razorpay docs:
//   - https://razorpay.com/docs/build/llm-docs/payments/payments/failure-analysis.md
//   - https://razorpay.com/docs/build/llm-docs/errors/payments/payment-methods-error-parameters.md
var hardDeclineCodes = map[string]bool{
	// Issuer bank hard blocks — docs say "customer must reach out to issuer bank"
	"card_declined":    true, // Docs: "Issuer Banks can decline the card due to multiple checks"
	"payment_declined": true, // Docs: "Issuer Bank or Gateway declined due to business/technical reasons"

	// Card-specific permanent blocks
	"card_number_invalid":               true,
	"card_disabled_for_online_payments": true,
	"expired_card":                      true, // Card past expiration date — no retry possible
	"lost_card":                         true, // Bank flagged as lost
	"stolen_card":                       true, // Bank flagged as stolen
	"do_not_honor":                      true, // ISO code 05 generic hard block
	"debit_instrument_blocked":          true,
	"pin_attempts_exceeded":             true, // Card locked after too many wrong PINs

	// Card not enrolled in 3DS — docs: "enroll card or use a different method"
	"card_not_enrolled": true,

	// Risk/fraud — retrying is harmful to merchant's risk profile
	"payment_risk_check_failed": true, // Docs: "retry with a different card or method"

	// Authentication failed permanently — requires a fresh checkout session
	"authentication_failed": true, // Docs: "customer has entered incorrect card details"
	"invalid_device":        true,

	// UPI / Account hard blocks
	"vpa_not_found":                      true, // UPI ID does not exist — must correct it first
	"beneficiary_account_does_not_exist": true,
	"upi_payment_cancelled":              true, // Deliberate cancel — do not auto-retry

	// Business / integration hard stops
	"input_validation_failed":               true, // Wrong integration parameters
	"international_transaction_not_allowed": true, // Account not enabled for intl payments
	"invalid_amount":                        true, // Invalid amount in request
	"invalid_currency":                      true, // Currency not enabled
	"mobile_number_invalid":                 true, // Docs: "check mobile mapped to UPI account"
	"duplicate_request":                     true, // Same idempotency key — do not retry
}

// RazorpayPaymentEvent is a targeted struct that maps the nested Razorpay webhook JSON structure.
// Using struct tags guarantees safe, type-checked parsing without fragile map assertions.
// Unknown fields in the payload are automatically ignored by json.Unmarshal.
type RazorpayPaymentEvent struct {
	Payload struct {
		Payment struct {
			Entity struct {
				ID          string `json:"id"`
				ErrorCode   string `json:"error_code"`
				ErrorReason string `json:"error_reason"`
			} `json:"entity"`
		} `json:"payment"`
	} `json:"payload"`
}

// ClassifyPayload takes the raw Razorpay webhook JSON, safely extracts the
// necessary fields via struct tags, and categorizes the decline as Soft, Hard,
// or Unknown. Unknown means the error_reason is not in either list and should
// be forwarded to an AI agent for a final decision.
func ClassifyPayload(rawJSON []byte) (*ClassificationResult, error) {
	var event RazorpayPaymentEvent

	if err := json.Unmarshal(rawJSON, &event); err != nil {
		return nil, err
	}

	entity := event.Payload.Payment.Entity

	result := &ClassificationResult{
		Type:        Unknown, // Default — AI agent will decide if neither list matches
		PaymentID:   entity.ID,
		ErrorCode:   entity.ErrorCode,
		ErrorReason: entity.ErrorReason,
	}

	switch {
	case softDeclineCodes[result.ErrorReason]:
		result.Type = Soft
	case hardDeclineCodes[result.ErrorReason]:
		result.Type = Hard
	}

	return result, nil
}
