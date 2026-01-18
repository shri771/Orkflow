package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const OllamaTimeout = 3 * time.Minute // Max time for generation

type OllamaClient struct {
	Endpoint string
	Model    string
}

func (o *OllamaClient) Generate(prompt string) (string, error) {
	payload := map[string]interface{}{
		"model":  o.Model,
		"prompt": prompt,
		"stream": false,
	}

	body, _ := json.Marshal(payload)
	url := o.Endpoint + "/api/generate"

	// Create request with timeout
	ctx, cancel := context.WithTimeout(context.Background(), OllamaTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("TIMEOUT: Ollama generation exceeded %v (try a faster model or shorter prompt)", OllamaTimeout)
		}
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama api error: %s", string(respBody))
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Response, nil
}
