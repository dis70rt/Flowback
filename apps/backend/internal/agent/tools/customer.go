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

type GetCustomerInput struct {
	RazorpayCustomerID string `json:"razorpay_customer_id"`
}

type CustomerOutput struct {
	RazorpayCustomerID string  `json:"razorpay_customer_id"`
	Name               string  `json:"name"`
	ValueTier          string  `json:"value_tier"`
	Tenure             string  `json:"tenure"`
	PreferredChannel   string  `json:"preferred_channel"`
	City               string  `json:"city"`
	State              string  `json:"state"`
	FailedPayments     int32   `json:"failed_payments"`
	ReliabilityScore   float64 `json:"reliability_score"`
}

func NewGetCustomerTool(queries *repo.Queries) (tool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "get_customer_profile",
			Description: "Fetches the customer's LTV, tenure, location, and preferred channel.",
		},
		func(ctx agent.Context, input GetCustomerInput) (*CustomerOutput, error) {
			slog.InfoContext(ctx, "tool called: get_customer_profile", "razorpay_customer_id", input.RazorpayCustomerID)
			param := sql.NullString{String: input.RazorpayCustomerID, Valid: true}
			profile, err := queries.GetCustomerProfile(ctx, param)
			if err != nil {
				slog.WarnContext(ctx, "customer profile not found", "razorpay_customer_id", input.RazorpayCustomerID, "error", err)
				return nil, fmt.Errorf("failed to fetch customer: %v", err)
			}

			out := &CustomerOutput{
				RazorpayCustomerID: profile.RazorpayCustomerID.String,
				Name:               profile.Name.String,
				ValueTier:          profile.ValueTier.String,
				Tenure:             profile.Tenure.String,
				PreferredChannel:   profile.PreferredChannel.String,
				City:               profile.City.String,
				State:              profile.State.String,
				FailedPayments:     profile.FailedPayments,
				ReliabilityScore:   profile.ReliabilityScore,
			}
			slog.InfoContext(ctx, "customer profile fetched",
				"name", out.Name,
				"tier", out.ValueTier,
				"preferred_channel", out.PreferredChannel,
				"failed_payments", out.FailedPayments,
				"reliability_score", out.ReliabilityScore,
				"city", out.City,
			)
			return out, nil
		},
	)
}
