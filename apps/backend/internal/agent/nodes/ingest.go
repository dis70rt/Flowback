package nodes

import (
	"encoding/json"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/workflow"
)

// RazorpayEvent provides a clean, strictly-typed way to parse the webhook
type RazorpayEvent struct {
	Payload struct {
		Payment struct {
			Entity struct {
				ID         string  `json:"id"`
				CustomerID string  `json:"customer_id"`
				Contact    string  `json:"contact"`
				Amount     float64 `json:"amount"`
			} `json:"entity"`
		} `json:"payment"`
		Subscription struct {
			Entity struct {
				ID string `json:"id"`
			} `json:"entity"`
		} `json:"subscription"`
	} `json:"payload"`
}

func NewIngestNode() *workflow.FunctionNode {
	return workflow.NewFunctionNode(
		"IngestNode",
		func(ctx agent.Context, webhookJSON string) (string, error) {
			var rzp RazorpayEvent
			if err := json.Unmarshal([]byte(webhookJSON), &rzp); err == nil {
				
				// Handle missing customer_id by falling back to contact
				customerID := rzp.Payload.Payment.Entity.CustomerID
				if customerID == "" {
					customerID = rzp.Payload.Payment.Entity.Contact
				}

				_ = ctx.State().Set("customer_id", customerID)
				_ = ctx.State().Set("payment_id", rzp.Payload.Payment.Entity.ID)
				_ = ctx.State().Set("amount", int64(rzp.Payload.Payment.Entity.Amount))
				_ = ctx.State().Set("subscription_id", rzp.Payload.Subscription.Entity.ID)
			}
			
			// Pass the raw JSON string forward to the Strategy Agent LLM
			return webhookJSON, nil 
		},
		workflow.NodeConfig{},
	)
}
