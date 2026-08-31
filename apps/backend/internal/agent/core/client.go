package core

import (
	"context"
	"fmt"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	openaimodel "google.golang.org/adk/v2/model/openaimodel"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/genai"
)

const (
	openRouterBaseURL = "https://openrouter.ai/api/v1"
	DefaultModel      = "z-ai/glm-5.3-flash"
	CopywriterModel   = "z-ai/glm-5.3-flash" // Highly recommended for creative writing / tone
)

// ModelRegistry holds our initialized ADK model clients.
type ModelRegistry struct {
	FastModel     model.LLM // For simple tasks: classification, diagnosis
	SmartModel    model.LLM // For complex reasoning: strategy, decision making
	CreativeModel model.LLM // For writing human-facing text: emails, SMS
}

// InitModels sets up the OpenRouter HTTP connections exactly once.
func InitModels(ctx context.Context, apiKey, defaultModelName string) (*ModelRegistry, error) {
	if defaultModelName == "" {
		defaultModelName = DefaultModel
	}

	fastModel, err := openaimodel.NewModel(ctx, defaultModelName, &openaimodel.ClientConfig{
		APIKey:  apiKey,
		BaseURL: openRouterBaseURL,
	})
	if err != nil {
		return nil, fmt.Errorf("registry: fast model: %w", err)
	}

	smartModel, err := openaimodel.NewModel(ctx, defaultModelName, &openaimodel.ClientConfig{
		APIKey:  apiKey,
		BaseURL: openRouterBaseURL,
	})
	if err != nil {
		return nil, fmt.Errorf("registry: smart model: %w", err)
	}

	creativeModel, err := openaimodel.NewModel(ctx, CopywriterModel, &openaimodel.ClientConfig{
		APIKey:  apiKey,
		BaseURL: openRouterBaseURL,
	})
	if err != nil {
		return nil, fmt.Errorf("registry: creative model: %w", err)
	}

	return &ModelRegistry{
		FastModel:     fastModel,
		SmartModel:    smartModel,
		CreativeModel: creativeModel,
	}, nil
}

type AgentConfig struct {
	Name                  string
	Description           string
	Instruction           string
	OutputKey             string
	OutputSchema          *genai.Schema
	GenerateContentConfig *genai.GenerateContentConfig
}

func (r *ModelRegistry) NewFastAgent(cfg AgentConfig, tools []tool.Tool) (agent.Agent, error) {
	return llmagent.New(llmagent.Config{
		Name:                  cfg.Name,
		Model:                 r.FastModel,
		Description:           cfg.Description,
		Instruction:           cfg.Instruction,
		OutputKey:             cfg.OutputKey,
		OutputSchema:          cfg.OutputSchema,
		GenerateContentConfig: cfg.GenerateContentConfig,
		Tools:                 tools,
	})
}

func (r *ModelRegistry) NewSmartAgent(cfg AgentConfig, tools []tool.Tool) (agent.Agent, error) {
	return llmagent.New(llmagent.Config{
		Name:                  cfg.Name,
		Model:                 r.SmartModel,
		Description:           cfg.Description,
		Instruction:           cfg.Instruction,
		OutputKey:             cfg.OutputKey,
		OutputSchema:          cfg.OutputSchema,
		GenerateContentConfig: cfg.GenerateContentConfig,
		Tools:                 tools,
	})
}

// NewCreativeAgent creates an LlmAgent backed by a model optimized for copywriting.
func (r *ModelRegistry) NewCreativeAgent(cfg AgentConfig, tools []tool.Tool) (agent.Agent, error) {
	return llmagent.New(llmagent.Config{
		Name:                  cfg.Name,
		Model:                 r.CreativeModel,
		Description:           cfg.Description,
		Instruction:           cfg.Instruction,
		OutputKey:             cfg.OutputKey,
		OutputSchema:          cfg.OutputSchema,
		GenerateContentConfig: cfg.GenerateContentConfig,
		Tools:                 tools,
	})
}

func ptr[T any](v T) *T { return &v }
