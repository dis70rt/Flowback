package razorpay

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	Secret string
}

func NewWebhookHandler(secret string) *WebhookHandler {
	return &WebhookHandler{Secret: secret}
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

	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	event, _ := payload["event"].(string)
	log.Printf("SUCCESS: SECURE EVENT RECEIVED: %s\n", event)

	// TODO: Send to Kafka

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
