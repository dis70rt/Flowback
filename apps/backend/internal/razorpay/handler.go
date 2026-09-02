package razorpay

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dis70rt/flowback/internal/repo"
	"github.com/gin-gonic/gin"
)

type TaskEnqueuer interface {
	EnqueueWebhook(event string, rawJSON []byte) error
}

type WebhookHandler struct {
	Secret   string
	Enqueuer TaskEnqueuer
	Queries  *repo.Queries
}

func NewWebhookHandler(secret string, enqueuer TaskEnqueuer, queries *repo.Queries) *WebhookHandler {
	return &WebhookHandler{
		Secret:   secret,
		Enqueuer: enqueuer,
		Queries:  queries,
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
		slog.Info("ERROR: Webhook signature verification failed! Possible malicious request.")
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

	slog.InfoContext(c.Request.Context(), "secure event received", "event", payload.Event)

	// LOG WEBHOOK TO DB (THIS FIRES THE SMART POSTGRES TRIGGER FOR SUCCESS EVENTS!)
	err = h.Queries.LogWebhookEvent(c.Request.Context(), repo.LogWebhookEventParams{
		RazorpayEventID: c.GetHeader("X-Razorpay-Event-Id"), // Optional tracking
		EventType:       payload.Event,
		Payload:         json.RawMessage(body),
		Signature:       signature,
	})
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "failed to log webhook to DB", "error", err)
	}

	// ENQUEUE FOR AI ONLY IF IT IS A FAILED PAYMENT
	if payload.Event == "payment.failed" || payload.Event == "subscription.charged.failed" {
		if err := h.Enqueuer.EnqueueWebhook(payload.Event, body); err != nil {
			slog.ErrorContext(c.Request.Context(), "failed to enqueue webhook", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
