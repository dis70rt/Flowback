package nodes

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/dis70rt/flowback/internal/repo"
	"github.com/google/uuid"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/workflow"
)

type RazorpayEvent struct {
	Payload struct {
		Payment struct {
			Entity struct {
				ID         string  `json:"id"`
				CustomerID string  `json:"customer_id"`
				Contact    string  `json:"contact"`
				Email      string  `json:"email"`
				Amount     float64 `json:"amount"`
			} `json:"entity"`
		} `json:"payment"`
		Subscription struct {
			Entity struct {
				ID         string `json:"id"`
				CustomerID string `json:"customer_id"`
			} `json:"entity"`
		} `json:"subscription"`
		Order struct {
			Entity struct {
				ID         string `json:"id"`
				CustomerID string `json:"customer_id"`
			} `json:"entity"`
		} `json:"order"`
	} `json:"payload"`
}

func NewIngestNode(queries *repo.Queries) *workflow.FunctionNode {
	return workflow.NewFunctionNode(
		"IngestNode",
		func(ctx agent.Context, webhookJSON string) (string, error) {
			var rzp RazorpayEvent
			if err := json.Unmarshal([]byte(webhookJSON), &rzp); err == nil {
				
				email := rzp.Payload.Payment.Entity.Email
				phone := rzp.Payload.Payment.Entity.Contact
				
				// 1. Try to find the internal UUID in our DB using Email or Phone!
				var internalUUID uuid.NullUUID
				if email != "" || phone != "" {
					dbID, err := queries.GetCustomerByEmailOrPhone(ctx, repo.GetCustomerByEmailOrPhoneParams{
						Email: sql.NullString{String: email, Valid: email != ""},
						Phone: sql.NullString{String: phone, Valid: phone != ""},
					})
					if err == nil {
						internalUUID = uuid.NullUUID{UUID: dbID, Valid: true}
					} else if err != sql.ErrNoRows {
						log.Printf("ERROR looking up customer by email/phone: %v", err)
					}
				}

				// We attach the verified internal DB UUID directly to the whiteboard
				_ = ctx.State().Set("internal_customer_uuid", internalUUID)

				// We still extract the raw Razorpay strings just in case we need them for metrics
				rzpCustID := rzp.Payload.Subscription.Entity.CustomerID
				if rzpCustID == "" { rzpCustID = rzp.Payload.Order.Entity.CustomerID }
				if rzpCustID == "" { rzpCustID = rzp.Payload.Payment.Entity.CustomerID }
				if rzpCustID == "" { rzpCustID = phone }
				if rzpCustID == "" { rzpCustID = email }

				_ = ctx.State().Set("customer_id", rzpCustID) // fallback string for display
				_ = ctx.State().Set("payment_id", rzp.Payload.Payment.Entity.ID)
				_ = ctx.State().Set("amount", int64(rzp.Payload.Payment.Entity.Amount))
				_ = ctx.State().Set("subscription_id", rzp.Payload.Subscription.Entity.ID)
			}
			
			return webhookJSON, nil 
		},
		workflow.NodeConfig{},
	)
}
