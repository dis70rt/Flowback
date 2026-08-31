package agent

import (
	"context"

	"github.com/dis70rt/flowback/internal/repo"
	"github.com/google/uuid"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/workflowagent"
	"google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/workflow"
	
	"github.com/dis70rt/flowback/internal/agent/core"
)

type RecoveryState struct {
	CaseID uuid.UUID
}

func BuildOrchestrator(
	ctx context.Context, 
	queries *repo.Queries,
	strategyAgent agent.Agent,
	copywriterAgent agent.Agent,
) (agent.Agent, error) {

	nodeStrategy, err := workflow.NewAgentNode(strategyAgent, workflow.NodeConfig{})
	if err != nil {
		return nil, err
	}
	
	nodeCopywriter, err := workflow.NewAgentNode(copywriterAgent, workflow.NodeConfig{})
	if err != nil {
		return nil, err
	}

	policyGuardrail := workflow.NewEmittingFunctionNode(
		"PolicyGuardrail",
		func(ctx agent.Context, strategyOut core.StrategyOutput, emit func(*session.Event) error) (any, error) {
			
			route := "copywriter"
			if strategyOut.Action == "silent_retry" {
				route = "execute"
			}

			ev := session.NewEvent(ctx, ctx.InvocationID())
			ev.Routes = []string{route}
			ev.Output = strategyOut 
			if err := emit(ev); err != nil {
				return nil, err
			}
			return nil, nil
		},
		workflow.NodeConfig{},
	)

	executionNode := workflow.NewFunctionNode(
		"ExecutionNode",
		func(ctx agent.Context, draft any) (string, error) {
			return "Success", nil
		},
		workflow.NodeConfig{},
	)

	edges := workflow.Concat(
		workflow.Chain(workflow.Start, nodeStrategy, policyGuardrail),
		[]workflow.Edge{
			{From: policyGuardrail, To: nodeCopywriter, Route: workflow.StringRoute("copywriter")},
			{From: policyGuardrail, To: executionNode, Route: workflow.StringRoute("execute")},
			{From: nodeCopywriter, To: executionNode},
		},
	)

	return workflowagent.New(workflowagent.Config{
		Name:        "orchestrator_agent",
		Description: "Manages the entire Razorpay webhook lifecycle.",
		Edges:       edges,
	})
}
