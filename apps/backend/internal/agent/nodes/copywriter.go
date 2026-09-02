package nodes

import (
	"github.com/dis70rt/flowback/internal/agent/core"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/genai"
)

// CopywriterOutput is the typed Go struct for the generated message.
type CopywriterOutput struct {
	Subject string `json:"subject"` // Empty if channel is SMS
	Body    string `json:"body"`    // The personalized message content
}

// copywriterSchema enforces the shape of the message generation.
var copywriterSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"subject": {Type: genai.TypeString},
		"body":    {Type: genai.TypeString},
	},
	Required: []string{"subject", "body"},
}

// newCopywriterAgent is package-private. Only NewAgents can call it.
// Uses the CreativeModel — tone, empathy, and clarity matter here.
func NewCopywriterAgent(registry *core.ModelRegistry) (agent.Agent, error) {
	cfg := core.AgentConfig{
		Name:        "copywriter_agent",
		Description: "Drafts a personalized recovery message (email/SMS) for the customer",
		Instruction: copywriterPrompt,

		// The orchestrator can read {message} to actually send it via execution tools
		OutputKey: "message",

		OutputSchema: copywriterSchema,

		// Higher temperature (0.6) allows for natural, empathetic, and varied human-like writing
		GenerateContentConfig: &genai.GenerateContentConfig{
			Temperature: core.Ptr[float32](0.6),
		},
	}

	return registry.NewCreativeAgent(cfg, nil)
}
