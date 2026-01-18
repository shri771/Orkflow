package agent

type LLMClient interface {
	Generate(prompt string) (string, error)
}

func NewLLMClient(provider string, model string, apiKey string, endpoint string) LLMClient {
	switch provider {
	case "anthropic":
		return &ClaudeClient{
			APIKey: apiKey,
			Model:  model,
		}
	case "openai":
		return &OpenAIClient{
			APIKey: apiKey,
			Model:  model,
		}
	case "gemini", "google":
		return &GeminiClient{
			APIKey: apiKey,
			Model:  model,
		}
	case "ollama":
		ep := endpoint
		if ep == "" {
			ep = "http://localhost:11434"
		}
		return &OllamaClient{
			Endpoint: ep,
			Model:    model,
		}
	default:
		// Use generic OpenAI-compatible client for any other provider
		// This auto-handles: groq, mistral, together, perplexity, openrouter, etc.
		if apiKey != "" {
			return NewGenericClient(provider, model, apiKey, endpoint)
		}
		// Fallback to Ollama if no API key
		return &OllamaClient{
			Endpoint: "http://localhost:11434",
			Model:    model,
		}
	}
}
