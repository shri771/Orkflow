package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Known OpenAI-compatible API endpoints
var knownEndpoints = map[string]string{
	"groq":       "https://api.groq.com/openai/v1/chat/completions",
	"mistral":    "https://api.mistral.ai/v1/chat/completions",
	"together":   "https://api.together.xyz/v1/chat/completions",
	"perplexity": "https://api.perplexity.ai/chat/completions",
	"openrouter": "https://openrouter.ai/api/v1/chat/completions",
	"deepseek":   "https://api.deepseek.com/v1/chat/completions",
	"fireworks":  "https://api.fireworks.ai/inference/v1/chat/completions",
}

// GenericClient is an OpenAI-compatible client that works with many providers
type GenericClient struct {
	APIKey   string
	Model    string
	Endpoint string
	Provider string
}

// NewGenericClient creates a client for any OpenAI-compatible API
func NewGenericClient(provider, model, apiKey, endpoint string) *GenericClient {
	// Use known endpoint if available, otherwise use provided endpoint
	ep := endpoint
	if ep == "" {
		if known, ok := knownEndpoints[strings.ToLower(provider)]; ok {
			ep = known
		} else {
			// Default to OpenAI format: https://api.{provider}.com/v1/chat/completions
			ep = fmt.Sprintf("https://api.%s.com/v1/chat/completions", strings.ToLower(provider))
		}
	}

	return &GenericClient{
		APIKey:   apiKey,
		Model:    model,
		Endpoint: ep,
		Provider: provider,
	}
}

func (g *GenericClient) Generate(prompt string) (string, error) {
	payload := map[string]interface{}{
		"model": g.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", g.Endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+g.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s connection error: %w", g.Provider, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// Check for common errors
		errStr := string(respBody)
		if strings.Contains(errStr, "invalid_api_key") || strings.Contains(errStr, "Unauthorized") {
			return "", fmt.Errorf("%s: invalid API key", g.Provider)
		}
		if strings.Contains(errStr, "rate_limit") || strings.Contains(errStr, "quota") {
			return "", fmt.Errorf("QUOTA_EXCEEDED[%s]: rate limit reached", g.Provider)
		}
		return "", fmt.Errorf("%s API error (%d): %s", g.Provider, resp.StatusCode, errStr)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("%s: failed to parse response: %w", g.Provider, err)
	}

	if result.Error.Message != "" {
		return "", fmt.Errorf("%s error: %s", g.Provider, result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from %s", g.Provider)
	}

	return result.Choices[0].Message.Content, nil
}
