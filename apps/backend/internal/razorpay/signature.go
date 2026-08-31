package razorpay

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
)

// VerifySignature verifies the Razorpay webhook signature using HMAC-SHA256.
// Uses constant-time comparison to prevent timing attacks.
func VerifySignature(body []byte, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	return subtle.ConstantTimeCompare([]byte(expectedSignature), []byte(signature)) == 1
}
