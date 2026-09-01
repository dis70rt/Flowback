package nodes

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/dis70rt/flowback/internal/repo"
	"github.com/google/uuid"
	"github.com/dis70rt/flowback/internal/agent/tools"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/workflow"
)


type EnrichedPayload struct {
	Webhook  string         `json:"webhook"`
	Customer *repo.Customer `json:"customer_profile"`
	NewsContext []string `json:"local_news_headlines,omitempty"`
}

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
				
				rzpCustID := rzp.Payload.Subscription.Entity.CustomerID
				if rzpCustID == "" { rzpCustID = rzp.Payload.Order.Entity.CustomerID }
				if rzpCustID == "" { rzpCustID = rzp.Payload.Payment.Entity.CustomerID }

				var internalUUID uuid.NullUUID

				if rzpCustID != "" {
					profile, err := queries.GetCustomerProfile(ctx, sql.NullString{String: rzpCustID, Valid: true})
					if err == nil {
						internalUUID = uuid.NullUUID{UUID: profile.ID, Valid: true}
					}
				}

				if !internalUUID.Valid && (email != "" || phone != "") {
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
				_ = ctx.State().Set("internal_customer_uuid", internalUUID)

				if rzpCustID == "" { rzpCustID = phone }
				if rzpCustID == "" { rzpCustID = email }

				log.Printf("Ingested Razorpay webhook: customer_id=%s, razorpay_customer_id=%s, payment_id=%s, subscription_id=%s, amount=%f", internalUUID.UUID.String(), rzpCustID, rzp.Payload.Payment.Entity.ID, rzp.Payload.Subscription.Entity.ID, rzp.Payload.Payment.Entity.Amount)

				_ = ctx.State().Set("customer_id", rzpCustID) // fallback string for display
				_ = ctx.State().Set("payment_id", rzp.Payload.Payment.Entity.ID)
				_ = ctx.State().Set("amount", int64(rzp.Payload.Payment.Entity.Amount))
				_ = ctx.State().Set("subscription_id", rzp.Payload.Subscription.Entity.ID)
			}
			
			// Bundle webhook and customer profile into a single JSON
			enriched := EnrichedPayload{Webhook: webhookJSON}
			
			internalUUIDVal, err := ctx.State().Get("internal_customer_uuid")
			if err == nil {
				if internalUUID, ok := internalUUIDVal.(uuid.NullUUID); ok && internalUUID.Valid {
					cust, err := queries.GetCustomerByID(ctx, internalUUID.UUID)
					if err == nil {
						enriched.Customer = &cust
						
						// Fetch news context automatically
						if cust.City.String != "" {
							headlines, _, err := tools.FetchLocalNews(cust.City.String, "bank outage internet storm flood")
							if err == nil {
								enriched.NewsContext = headlines
							} else {
								log.Printf("Warning: failed to fetch news context for %s: %v", cust.City.String, err)
							}
						}
					}
				}
			}
			
			b, _ := json.MarshalIndent(enriched, "", "  ")
			return string(b), nil 
		},
		workflow.NodeConfig{},
	)
}
