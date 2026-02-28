package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type NvidiaClient struct {
	apiKey  string
	client  *http.Client
	baseURL string
}

type ChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewNvidiaClient initializes a new NvidiaClient with API key and default HTTP client
func NewNvidiaClient() *NvidiaClient {
	return &NvidiaClient{
		apiKey: os.Getenv("NVIDIA_API_KEY"),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://integrate.api.nvidia.com/v1/chat/completions",
	}
}

// ChatCompletion sends a prompt to NVIDIA AI and returns the response
func (n *NvidiaClient) ChatCompletion(ctx context.Context, prompt string) (string, error) {
	payload := map[string]interface{}{
		"model": "nvidia/nemotron-3-nano-30b-a3b",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 1,
		"top_p":       1,
		"max_tokens":  512,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", n.baseURL, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+n.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(responseData))
	}

	var aiResp ChatCompletionResponse
	if err := json.Unmarshal(responseData, &aiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal AI response: %w", err)
	}

	if len(aiResp.Choices) == 0 {
		return "", errors.New("no AI response returned")
	}

	// Trim any extra whitespace from AI response
	return strings.TrimSpace(aiResp.Choices[0].Message.Content), nil
}