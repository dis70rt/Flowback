package decline

// Type represents whether a payment failure is recoverable or not
type Type string

const (
	// Soft means the failure is transient - the card/account itself is fine.
	// Action: auto-retry or notify the customer to try again soon.
	Soft Type = "soft"

	// Hard means the failure is permanent - the card/account is fundamentally blocked.
	// Action: require the customer to provide a new payment method.
	Hard Type = "hard"

	// Unknown is returned when the error code is not recognized.
	Unknown Type = "unknown"
)

// softDeclineCodes are transient failures where retrying has a realistic chance of success.
// Source: Razorpay test card and UPI error code documentation.
var softDeclineCodes = map[string]bool{
	// Infra / Bank side transient issues - retry automatically
	"payment_timed_out":       true,
	"bank_technical_error":    true,
	"gateway_technical_error": true,
	"upi_app_not_available":   true,

	// Customer fixable issues - notify and let them retry
	"insufficient_fund":                    true,
	"transaction_limit_exceeded":           true,
	"transaction_frequency_limit_exceeded": true,
	"incorrect_pin":                        true,
	"pin_not_set":                          true,
}

// hardDeclineCodes are permanent failures. Auto-retry will never succeed.
// The customer must provide a new payment method.
var hardDeclineCodes = map[string]bool{
	"card_declined":                      true,
	"payment_declined":                   true,
	"card_disabled_for_online_payments":  true,
	"card_number_invalid":                true,
	"authentication_failed":              true,
	"debit_instrument_blocked":           true,
	"payment_risk_check_failed":          true,
	"invalid_device":                     true,
	"beneficiary_account_does_not_exist": true,
	"pin_attempts_exceeded":              true,
	"duplicate_request":                  true,
	"payment_cancelled":                  true,
}

// Classify takes a Razorpay error_reason string and deterministically returns
// the decline type. This is a pure lookup function with no external dependencies.
func Classify(errorReason string) Type {
	if softDeclineCodes[errorReason] {
		return Soft
	}
	if hardDeclineCodes[errorReason] {
		return Hard
	}
	return Unknown
}
