package nodes

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dis70rt/flowback/internal/agent/core"
	"google.golang.org/adk/v2/agent"
)

//go:embed prompts/hinglish_voice.md
var hinglishVoicePrompt string

func NewVoiceAgent(registry *core.ModelRegistry) (agent.Agent, error) {
	cfg := core.AgentConfig{
		Name:         "voice_script_agent",
		Description:  "Drafts a conversational Hinglish voice script for the customer.",
		Instruction:  hinglishVoicePrompt,
		OutputKey:    "message",
		OutputSchema: copywriterSchema, // reuse copywriter schema so it returns {"message": "..."}
	}

	return registry.NewCreativeAgent(cfg, nil)
}

// GenerateVoiceAudio is a helper for the Hackathon Demo.
// It directly calls OpenRouter's TTS endpoint to get the base64 audio bytes.
func GenerateVoiceAudio(script string, apiKey string) (string, error) {
	url := "https://openrouter.ai/api/v1/audio/speech"

	payload := map[string]interface{}{
		"model":          "google/gemini-3.1-flash-tts-preview",
		"input":          script,
		"responseFormat": "wav",
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openrouter TTS api failed: %d - %s", resp.StatusCode, string(b))
	}

	audioBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return "data:audio/wav;base64," + base64.StdEncoding.EncodeToString(audioBytes), nil
}
