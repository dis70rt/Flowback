package nodes

import (
	"google.golang.org/adk/v2/agent"
	"google.golang.org/genai"
	"github.com/dis70rt/flowback/internal/agent/core"
)

type DiagnoseOutput struct {
	Category      string  `json:"category"`       // soft_decline | hard_decline | unknown
	Reasoning     string  `json:"reasoning"`      // one-sentence explanation
	IsRecoverable bool    `json:"is_recoverable"` // drives the workflow router
	Confidence    float32 `json:"confidence"`     // 0.0 – 1.0
}

var diagnoseSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"category": {
			Type: genai.TypeString,
			Enum: []string{"soft_decline", "hard_decline", "unknown"},
		},
		"reasoning":      {Type: genai.TypeString},
		"is_recoverable": {Type: genai.TypeBoolean},
		"confidence": {
			Type: genai.TypeNumber,
		},
	},
	Required: []string{"category", "reasoning", "is_recoverable", "confidence"},
}

func NewDiagnoseAgent(registry *core.ModelRegistry) (agent.Agent, error) {
	cfg := core.AgentConfig{
		Name:        "diagnose_agent",
		Description: "Classifies the root cause of a failed subscription payment",
		Instruction: diagnosePrompt, // embedded from prompts/diagnose.md at compile time

		OutputKey: "diagnosis",

		OutputSchema: diagnoseSchema,

		GenerateContentConfig: &genai.GenerateContentConfig{
			Temperature: core.Ptr[float32](0.1),
		},
	}

	return registry.NewFastAgent(cfg, nil)
}
