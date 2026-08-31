package nodes

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/dis70rt/flowback/internal/pubsub"
	"github.com/dis70rt/flowback/internal/repo"
	"github.com/google/uuid"
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
	default:
		return repo.ActionTypeESCALATETOHUMAN
	}
}

func NewExecutionNode(queries *repo.Queries, bus pubsub.Publisher) *workflow.FunctionNode {
	return workflow.NewFunctionNode(
		"ExecutionNode",
		func(ctx agent.Context, draft any) (string, error) {
			
			channelVal, _ := ctx.State().Get("channel")
			channel, _ := channelVal.(string)
			
			reasoningVal, _ := ctx.State().Get("reasoning")
			reasoning, _ := reasoningVal.(string)

			customerIDVal, _ := ctx.State().Get("customer_id")
			customerID, _ := customerIDVal.(string)

			subscriptionIDVal, _ := ctx.State().Get("subscription_id")
			subscriptionID, _ := subscriptionIDVal.(string)

			paymentIDVal, _ := ctx.State().Get("payment_id")
			paymentID, _ := paymentIDVal.(string)

			amountVal, _ := ctx.State().Get("amount")
			amount, _ := amountVal.(int64)

			// DB Operation 1: Create Recovery Case
			caseID, err := queries.CreateRecoveryCase(ctx, repo.CreateRecoveryCaseParams{
				CustomerID:        uuid.NullUUID{}, 
				SubscriptionID:    subscriptionID,
				PaymentID:         sql.NullString{String: paymentID, Valid: true},
				AmountAtRisk:      amount,
				Currency:          "INR",
			})
			if err != nil {
				log.Printf("ERROR creating recovery case: %v", err)
			}

			// DB Operation 2: Create Recovery Action
			draftBytes, _ := json.Marshal(draft)
			_, err = queries.CreateRecoveryAction(ctx, repo.CreateRecoveryActionParams{
				RecoveryCaseID: caseID,
				IdempotencyKey: ctx.InvocationID(),
				ActionType:     mapActionType(channel),
				Channel:        sql.NullString{String: channel, Valid: true},
				AiReasoning:    sql.NullString{String: reasoning, Valid: true},
				DraftBody:      sql.NullString{String: string(draftBytes), Valid: true},
				Status:         repo.ActionStatusPENDING,
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
			if err := bus.Publish(ctx, "dashboard_updates", payload); err != nil {
				return "", err
			}

			return "Success", nil
		},
		workflow.NodeConfig{},
	)
}
