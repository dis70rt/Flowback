package tools

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/dis70rt/flowback/internal/repo"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

type GetCommHistoryInput struct {
	RazorpayCustomerID string `json:"razorpay_customer_id"`
}

type CommHistoryOutput struct {
	Channel   string `json:"channel"`
	Status    string `json:"status"`
	SentAt    string `json:"sent_at"`
	OpenedAt  string `json:"opened_at,omitempty"`
	ClickedAt string `json:"clicked_at,omitempty"`
}

func NewGetCommunicationHistoryTool(queries *repo.Queries) (tool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "get_communication_history",
			Description: "Fetches past emails/messages to see if the user opened them or ignored them.",
		},
		func(ctx agent.Context, input GetCommHistoryInput) (*[]CommHistoryOutput, error) {
			slog.InfoContext(ctx, "tool called: get_communication_history", "razorpay_customer_id", input.RazorpayCustomerID)
			param := sql.NullString{String: input.RazorpayCustomerID, Valid: true}
			rows, err := queries.GetCommunicationHistory(ctx, param)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch history: %v", err)
			}

			var out []CommHistoryOutput
			for _, r := range rows {
				item := CommHistoryOutput{
					Channel: r.Channel,
					Status:  r.Status,
					SentAt:  r.SentAt.String(),
				}
				if r.OpenedAt.Valid {
					item.OpenedAt = r.OpenedAt.Time.String()
				}
				if r.ClickedAt.Valid {
					item.ClickedAt = r.ClickedAt.Time.String()
				}
				out = append(out, item)
			}

			slog.InfoContext(ctx, "communication history fetched", "razorpay_customer_id", input.RazorpayCustomerID, "record_count", len(out))
			return &out, nil
		},
	)
}
