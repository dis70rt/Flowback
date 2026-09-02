package nodes

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
		OutputSchema: copywriterSchema,
	}

	return registry.NewCreativeAgent(cfg, nil)
}

func prependWavHeader(pcmData []byte, sampleRate, channels, bitsPerSample int) []byte {
	header := make([]byte, 44)
	dataLen := len(pcmData)
	fileLen := dataLen + 36

	copy(header[0:4], []byte("RIFF"))
	binary.LittleEndian.PutUint32(header[4:8], uint32(fileLen))
	copy(header[8:12], []byte("WAVE"))

	copy(header[12:16], []byte("fmt "))
	binary.LittleEndian.PutUint32(header[16:20], 16)
	binary.LittleEndian.PutUint16(header[20:22], 1)
	binary.LittleEndian.PutUint16(header[22:24], uint16(channels))
	binary.LittleEndian.PutUint32(header[24:28], uint32(sampleRate))

	byteRate := sampleRate * channels * (bitsPerSample / 8)
	binary.LittleEndian.PutUint32(header[28:32], uint32(byteRate))

	blockAlign := channels * (bitsPerSample / 8)
	binary.LittleEndian.PutUint16(header[32:34], uint16(blockAlign))
	binary.LittleEndian.PutUint16(header[34:36], uint16(bitsPerSample))

	copy(header[36:40], []byte("data"))
	binary.LittleEndian.PutUint32(header[40:44], uint32(dataLen))

	return append(header, pcmData...)
}

func GenerateVoiceAudio(script string, apiKey string) (string, error) {
	url := "https://openrouter.ai/api/v1/audio/speech"

	payload := map[string]interface{}{
		"model":          "google/gemini-3.1-flash-tts-preview",
		"input":          script,
		"voice":          "Zephyr",
		"response_format": "pcm",
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

	pcmBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Gemini TTS default PCM specs: 24000Hz, 1 channel, 16-bit
	wavBytes := prependWavHeader(pcmBytes, 24000, 1, 16)

	os.MkdirAll("./audio", 0755)
	filename := fmt.Sprintf("%d.wav", time.Now().UnixNano())
	os.WriteFile("./audio/" + filename, wavBytes, 0644)
	return "http://localhost:8080/audio/" + filename, nil
}
