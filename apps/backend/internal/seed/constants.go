package seed

// ==========================================
// UPI Constants
// ==========================================

const (
	UPIVPASuccess = "success@razorpay"
	UPIVPAFailure = "failure@razorpay"
)

// UPI Error Amounts (in Paise) to trigger specific BAD_REQUEST_ERRORs
const (
	UPIErrorIncorrectPIN               = 204
	UPIErrorPINNotSet                  = 205
	UPIErrorPINAttemptsExceeded        = 206
	UPIErrorTransactionLimitExceeded   = 208 // or 209
	UPIErrorFreqLimitExceeded          = 210
	UPIErrorDebitInstrumentBlocked     = 212
	UPIErrorPaymentDeclined            = 304
	UPIErrorInvalidDevice              = 407
)

// UPI Error Amounts (in Paise) to trigger specific GATEWAY_ERRORs
const (
	UPIGatewayBankTechnicalError       = 104 // or 106
	UPIGatewayPaymentTimedOut          = 105
	UPIGatewayAppNotAvailable          = 107
	UPIGatewayBeneficiaryBlocked       = 211
	UPIGatewayBeneficiaryDoesNotExist  = 213
	UPIGatewayRiskCheckFailed          = 404 // or 405
	UPIGatewayDuplicateRequest         = 406
)

// ==========================================
// Card Constants
// ==========================================

// Standard Success Test Cards
const (
	CardSuccessVisa       = "4100280000001007"
	CardSuccessMastercard = "5555510000081006"
	CardSuccessRuPay      = "6527658900001005"
	CardSuccessAmex       = "340256000401007"
)

// Subscription Test Cards
const (
	CardSubscriptionDomesticVisa         = "4718609108204366"
	CardSubscriptionInternationalMaster  = "5104015555555558"
)

// Card Error Numbers to trigger specific BAD_REQUEST_ERRORs
const (
	CardErrorPaymentTimedOut           = "4100280000090000" // Visa
	CardErrorInsufficientFund          = "4100280000080001" // Visa
	CardErrorPaymentCancelled          = "4100280000070002" // Visa
	CardErrorCardDeclined              = "4100280000060003" // Visa
	CardErrorCardDisabledForOnline     = "4100280000030006" // Visa
	CardErrorCardNumberInvalid         = "4100280000010008" // Visa
)

// Card Error Numbers to trigger specific GATEWAY_ERRORs
const (
	CardGatewayTechnicalError          = "4100280000020007" // Visa
	CardGatewayAuthenticationFailed    = "4100280000000009" // Visa
)

// Shared Test Card Details
const (
	TestCardExpiryMonth = "12"
	TestCardExpiryYear  = "30" // Future year (current is 2026)
	TestCardCVV         = "123"
)
