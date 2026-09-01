package agent

import (
	"context"

	"github.com/dis70rt/flowback/internal/repo"
	"github.com/dis70rt/flowback/internal/pubsub"
	"github.com/google/uuid"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/workflowagent"
	"google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/workflow"
	
	"github.com/dis70rt/flowback/internal/agent/core"
	"github.com/dis70rt/flowback/internal/agent/nodes"
)

type RecoveryState struct {
	CaseID uuid.UUID
}

func BuildOrchestrator(
	ctx context.Context, 
	queries *repo.Queries,
	strategyAgent agent.Agent,
	copywriterAgent agent.Agent,
	voiceAgent agent.Agent,
	bus pubsub.Publisher,
	openRouterAPIKey string,
) (agent.Agent, error) {

	ingestNode := nodes.NewIngestNode(queries)
	execDirect := nodes.NewDirectExecutionNode(queries, bus)
	execCopywriter := nodes.NewCopywriterExecutionNode(queries, bus)
	execVoice := nodes.NewVoiceExecutionNode(queries, bus, openRouterAPIKey)

	nodeStrategy, err := workflow.NewAgentNode(strategyAgent, workflow.NodeConfig{})
	if err != nil {
		return nil, err
	}
	
	nodeCopywriter, err := workflow.NewAgentNode(copywriterAgent, workflow.NodeConfig{})
	if err != nil {
		return nil, err
	}

	nodeVoice, err := workflow.NewAgentNode(voiceAgent, workflow.NodeConfig{})
	if err != nil {
		return nil, err
	}

	policyGuardrail := workflow.NewEmittingFunctionNode(
		"PolicyGuardrail",
		func(ctx agent.Context, strategyOut core.StrategyOutput, emit func(*session.Event) error) (any, error) {
			route := "copywriter"
			if strategyOut.Action == "silent_retry" {
				route = "execute"
			} else if strategyOut.Action == "send_call" {
				route = "voice"
			}
			
			_ = ctx.State().Set("channel", strategyOut.Action)
			_ = ctx.State().Set("reasoning", strategyOut.Reasoning)
			_ = ctx.State().Set("discount_percentage", strategyOut.DiscountPercentage)

			ev := session.NewEvent(ctx, ctx.InvocationID())
			ev.Routes = []string{route}
			
			if route == "copywriter" || route == "voice" {
				profileStr := "Name: Unknown, Tier: Basic"
				uid, err := ctx.State().Get("internal_customer_uuid")
				if err == nil {
					if u, ok := uid.(uuid.NullUUID); ok && u.Valid {
						cust, err := queries.GetCustomerByID(ctx, u.UUID)
						if err == nil {
							profileStr = "Name: " + cust.Name.String + ", Tier: " + cust.ValueTier.String + ", City: " + cust.City.String
						}
					}
				}
				
				ev.Output = map[string]string{
					"Action": strategyOut.Action,
					"StrategyReasoning": strategyOut.Reasoning,
					"CustomerProfile": profileStr,
				}
			} else {
				ev.Output = strategyOut 
			}
			
			if err := emit(ev); err != nil {
				return nil, err
			}
			return nil, nil
		},
		workflow.NodeConfig{},
	)

	edges := workflow.Concat(
		workflow.Chain(workflow.Start, ingestNode, nodeStrategy, policyGuardrail),
		[]workflow.Edge{
			{From: policyGuardrail, To: nodeCopywriter, Route: workflow.StringRoute("copywriter")},
			{From: policyGuardrail, To: nodeVoice, Route: workflow.StringRoute("voice")},
			{From: policyGuardrail, To: execDirect, Route: workflow.StringRoute("execute")},
			{From: nodeCopywriter, To: execCopywriter},
			{From: nodeVoice, To: execVoice},
		},
	)

	return workflowagent.New(workflowagent.Config{
		Name:        "orchestrator_agent",
		Description: "Manages the entire Razorpay webhook lifecycle.",
		Edges:       edges,
	})
}
