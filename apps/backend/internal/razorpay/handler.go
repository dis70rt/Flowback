package razorpay

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskEnqueuer interface {
	EnqueueWebhook(event string, rawJSON []byte) error
}

type WebhookHandler struct {
	Secret   string
	Enqueuer TaskEnqueuer
}

func NewWebhookHandler(secret string, enqueuer TaskEnqueuer) *WebhookHandler {
	return &WebhookHandler{
		Secret:   secret,
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
		Event string `json:"event"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	log.Printf("SUCCESS: SECURE EVENT RECEIVED: %s\n", payload.Event)

	if err := h.Enqueuer.EnqueueWebhook(payload.Event, body); err != nil {
		log.Printf("ERROR: Failed to enqueue webhook: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
