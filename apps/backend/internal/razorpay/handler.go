package razorpay

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dis70rt/flowback/internal/events"
)

type TaskEnqueuer interface {
	EnqueuePaymentFailed(payload events.PaymentFailedPayload) error
	EnqueueSubscriptionHalted(payload events.SubscriptionHaltedPayload) error
}

type WebhookHandler struct {
	Secret    string
	Enqueuer TaskEnqueuer
}

func NewWebhookHandler(secret string, enqueuer TaskEnqueuer) *WebhookHandler {
	return &WebhookHandler{
		Secret:    secret,
		Enqueuer: enqueuer,
	}
}

func (h *WebhookHandler) Handle(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read body"})
		return
	}

	signature := c.GetHeader("X-Razorpay-Signature")

	if !VerifySignature(body, signature, h.Secret) {
		log.Println("ERROR: Webhook signature verification failed! Possible malicious request.")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	var payload struct {
		Event   string                 `json:"event"`
		Payload map[string]interface{} `json:"payload"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	log.Printf("SUCCESS: SECURE EVENT RECEIVED: %s\n", payload.Event)

	// Route events cleanly using our Enqueuer abstraction
	switch payload.Event {
	case "payment.failed":
		var paymentId, errCode, errDesc, customerId string
		if paymentData, ok := payload.Payload["payment"].(map[string]interface{}); ok {
			if entity, ok := paymentData["entity"].(map[string]interface{}); ok {
				paymentId, _ = entity["id"].(string)
				errCode, _ = entity["error_code"].(string)
				errDesc, _ = entity["error_description"].(string)
				customerId, _ = entity["customer_id"].(string)
			}
		}

		_ = h.Enqueuer.EnqueuePaymentFailed(events.PaymentFailedPayload{
			PaymentID:  paymentId,
			ErrorCode:  errCode,
			ErrorDesc:  errDesc,
			CustomerID: customerId,
		})

	case "subscription.halted":
		var subId, customerId string
		if subData, ok := payload.Payload["subscription"].(map[string]interface{}); ok {
			if entity, ok := subData["entity"].(map[string]interface{}); ok {
				subId, _ = entity["id"].(string)
				customerId, _ = entity["customer_id"].(string)
			}
		}

		_ = h.Enqueuer.EnqueueSubscriptionHalted(events.SubscriptionHaltedPayload{
			SubscriptionID: subId,
			CustomerID:     customerId,
		})
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
