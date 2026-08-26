package events

const (
	TopicWebhookReceived = "webhook:received"
)

type WebhookPayload struct {
	Event   string `json:"event"`
	RawJSON []byte `json:"raw_json"`
}
