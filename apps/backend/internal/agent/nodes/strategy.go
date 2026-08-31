package nodes

import (
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/genai"
	"github.com/dis70rt/flowback/internal/agent/core"
	"github.com/dis70rt/flowback/internal/repo"
	"github.com/dis70rt/flowback/internal/agent/tools"
)


// strategySchema is the genai.Schema enforced on the LLM's response.
var strategySchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"action": {
			Type: genai.TypeString,
			Enum: []string{"silent_retry", "send_email", "send_sms", "create_payment_link"},
		},
		"delay_hours": {Type: genai.TypeInteger},
		"reasoning":   {Type: genai.TypeString},
		"confidence": {
			Type: genai.TypeNumber,
		},
	},
	Required: []string{"action", "delay_hours", "reasoning", "confidence"},
}

func NewStrategyAgent(registry *core.ModelRegistry, queries *repo.Queries) (agent.Agent, error) {
	cfg := core.AgentConfig{
		Name:        "strategy_agent",
		Description: "Recommends the optimal recovery action for a soft-declined payment",
		Instruction: strategyPrompt, // embedded from prompts/strategy.md at compile time

		OutputKey: "strategy",

		OutputSchema: strategySchema,

		GenerateContentConfig: &genai.GenerateContentConfig{
			Temperature: core.Ptr[float32](0.3),
		},
	}

	getCustomer, err := tools.NewGetCustomerTool(queries)
	if err != nil {
		return nil, err
	}

	getCommHistory, err := tools.NewGetCommunicationHistoryTool(queries)
	if err != nil {
		return nil, err
	}

	searchNews, err := tools.NewSearchLocalNewsTool()
	if err != nil {
		return nil, err
	}

	return registry.NewSmartAgent(cfg, []tool.Tool{getCustomer, getCommHistory, searchNews})
}
