package nodes

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/dis70rt/flowback/internal/pubsub"
	"github.com/dis70rt/flowback/internal/repo"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/workflow"
)

func mapActionType(s string) repo.ActionType {
	switch s {
	case "silent_retry":
		return repo.ActionTypeSILENTRETRY
	case "create_payment_link":
		return repo.ActionTypeSENDPAYMENTLINK
	case "send_email":
		return repo.ActionTypeSENDEMAIL
	case "send_sms":
		return repo.ActionTypeSENDSMS
	case "send_whatsapp":
		return repo.ActionTypeSENDWHATSAPP
	case "send_call":
		return repo.ActionTypeSENDCALL
	default:
		return repo.ActionTypeESCALATETOHUMAN
	}
}

// saveAndPublish extracts state variables, commits the action to the database, and fires an SSE event.
func saveAndPublish(ctx agent.Context, draft any, audioURL string, queries *repo.Queries, bus pubsub.Publisher) (string, error) {
	channelVal, _ := ctx.State().Get("channel")
	channel, _ := channelVal.(string)

	reasoningVal, _ := ctx.State().Get("reasoning")
	reasoning, _ := reasoningVal.(string)

	discountVal, _ := ctx.State().Get("discount_percentage")
	discount, _ := discountVal.(int)

	customerIDVal, _ := ctx.State().Get("customer_id")
	customerID, _ := customerIDVal.(string)

	subscriptionIDVal, _ := ctx.State().Get("subscription_id")
	subscriptionID, _ := subscriptionIDVal.(string)

	paymentIDVal, _ := ctx.State().Get("payment_id")
	paymentID, _ := paymentIDVal.(string)

	amountVal, _ := ctx.State().Get("amount")
	amount, _ := amountVal.(int64)

	internalUUIDVal, _ := ctx.State().Get("internal_customer_uuid")
	internalUUID, _ := internalUUIDVal.(uuid.NullUUID)

	// DB Operation 1: Create Recovery Case
	caseID, err := queries.CreateRecoveryCase(ctx, repo.CreateRecoveryCaseParams{
		CustomerID:     internalUUID,
		SubscriptionID: subscriptionID,
		PaymentID:      sql.NullString{String: paymentID, Valid: true},
		AmountAtRisk:   amount,
		Currency:       "INR",
	})
	if err != nil {
		log.Printf("ERROR creating recovery case: %v", err)
	}

	// DB Operation 2: Create Recovery Action
	var draftBytes []byte
	if draft != nil {
		draftBytes, _ = json.Marshal(draft)
	}

	_, err = queries.CreateRecoveryAction(ctx, repo.CreateRecoveryActionParams{
		RecoveryCaseID:     caseID,
		IdempotencyKey:     ctx.InvocationID(),
		ActionType:         mapActionType(channel),
		Channel:            sql.NullString{String: channel, Valid: true},
		AiReasoning:        sql.NullString{String: reasoning, Valid: true},
		DiscountPercentage: sql.NullInt32{Int32: int32(discount), Valid: discount > 0},
		DraftPayload:       pqtype.NullRawMessage{RawMessage: draftBytes, Valid: len(draftBytes) > 0},
		Status:             repo.ActionStatusPENDING,
	})
	if err != nil {
		log.Printf("ERROR saving recovery action: %v", err)
	}

	// Publish SSE Event
	payload := map[string]any{
		"event":       "draft_ready",
		"case_id":     caseID.String(),
		"customer_id": customerID,
		"channel":     channel,
		"status":      "PENDING_APPROVAL",
	}
	if audioURL != "" {
		payload["audio_url"] = audioURL
	}
	
	if err := bus.Publish(ctx, "dashboard_updates", payload); err != nil {
		return "", err
	}

	return "Success", nil
}

// NewCopywriterExecutionNode handles storing text-based messages (Email, SMS, WhatsApp).
func NewCopywriterExecutionNode(queries *repo.Queries, bus pubsub.Publisher) *workflow.FunctionNode {
	return workflow.NewFunctionNode(
		"ExecCopywriter",
		func(ctx agent.Context, draft any) (string, error) {
			return saveAndPublish(ctx, draft, "", queries, bus)
		},
		workflow.NodeConfig{},
	)
}

// NewDirectExecutionNode handles automated actions with no draft (e.g. silent retry).
func NewDirectExecutionNode(queries *repo.Queries, bus pubsub.Publisher) *workflow.FunctionNode {
	return workflow.NewFunctionNode(
		"ExecDirect",
		func(ctx agent.Context, input any) (string, error) {
			return saveAndPublish(ctx, nil, "", queries, bus)
		},
		workflow.NodeConfig{},
	)
}

// NewVoiceExecutionNode handles TTS generation before saving.
func NewVoiceExecutionNode(queries *repo.Queries, bus pubsub.Publisher, openRouterAPIKey string) *workflow.FunctionNode {
	return workflow.NewFunctionNode(
		"ExecVoice",
		func(ctx agent.Context, draft any) (string, error) {
			draftStr := ""
			if m, ok := draft.(map[string]any); ok {
				if msg, ok := m["message"].(string); ok {
					draftStr = msg
				}
			}
			if draftStr == "" {
				if s, ok := draft.(string); ok {
					draftStr = s
				} else {
					b, _ := json.Marshal(draft)
					draftStr = string(b)
				}
			}

			log.Printf("Synthesizing voice audio using OpenRouter TTS...")
			audioURL := ""
			url, err := GenerateVoiceAudio(draftStr, openRouterAPIKey)
			if err != nil {
				log.Printf("ERROR generating audio: %v", err)
			} else {
				audioURL = url
				log.Printf("Audio Base64 successfully generated.")

				// Inject audioURL into draft so it saves to DB
				if m, ok := draft.(map[string]any); ok {
					m["audio_url"] = audioURL
					draft = m
				} else {
					draft = map[string]any{
						"message":   draftStr,
						"audio_url": audioURL,
					}
				}
			}

			return saveAndPublish(ctx, draft, audioURL, queries, bus)
		},
		workflow.NodeConfig{},
	)
}
