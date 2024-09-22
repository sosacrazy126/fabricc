package lmstudio

import (
	"context"
	"fmt"
	"net/http"
	"encoding/json"
	"bytes"
	"github.com/danielmiessler/fabric/common"
)

// NewClient creates a new LM Studio client with default configuration.
func NewClient() *Client {
	return NewClientCompatible("LM Studio", "http://localhost:1234/v1", nil)
}

// NewClientCompatible creates a new LM Studio client with custom configuration.
func NewClientCompatible(vendorName string, defaultBaseUrl string, configureCustom func() error) *Client {
	client := &Client{}
	if configureCustom == nil {
		configureCustom = client.configure
	}
	client.Configurable = &common.Configurable{
		Label:           vendorName,
		EnvNamePrefix:   common.BuildEnvVariablePrefix(vendorName),
		ConfigureCustom: configureCustom,
	}
	client.ApiBaseURL = client.AddSetupQuestion("API Base URL", false)
	client.ApiBaseURL.Value = defaultBaseUrl
	return client
}

// Client represents the LM Studio client.
type Client struct {
	*common.Configurable
	ApiBaseURL *common.SetupQuestion
	HttpClient *http.Client
}

// configure sets up the HTTP client.
func (c *Client) configure() error {
	c.HttpClient = &http.Client{}
	return nil
}

// Configure sets up the client configuration.
func (c *Client) Configure() error {
	return c.ConfigureCustom()
}

// ListModels returns a list of available models.
func (c *Client) ListModels() ([]string, error) {
	url := fmt.Sprintf("%s/models", c.ApiBaseURL.Value)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	
	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	models := make([]string, len(result.Data))
	for i, model := range result.Data {
		models[i] = model.ID
	}
	
	return models, nil
}

// SendStream sends a stream of messages (not implemented for LM Studio).
func (c *Client) SendStream(msgs []*common.Message, opts *common.ChatOptions, channel chan string) error {
	return fmt.Errorf("streaming is not currently supported for LM Studio")
}

// Send sends a single message and returns the response.
func (c *Client) Send(ctx context.Context, msgs []*common.Message, opts *common.ChatOptions) (string, error) {
	url := fmt.Sprintf("%s/chat/completions", c.ApiBaseURL.Value)
	
	payload := map[string]interface{}{
		"messages": msgs,
		"model":    opts.Model,
		// Add other options from opts if supported by LM Studio
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("invalid response format: missing or empty choices")
	}
	
	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format: missing message in first choice")
	}
	
	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response format: missing or non-string content in message")
	}
	
	return content, nil
}

// Complete sends a completion request and returns the response.
func (c *Client) Complete(ctx context.Context, prompt string, opts *common.ChatOptions) (string, error) {
	url := fmt.Sprintf("%s/completions", c.ApiBaseURL.Value)
	
	payload := map[string]interface{}{
		"prompt": prompt,
		"model":  opts.Model,
		// Add other options from opts if supported by LM Studio
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("invalid response format: missing or empty choices")
	}
	
	text, ok := choices[0].(map[string]interface{})["text"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response format: missing or non-string text in first choice")
	}
	
	return text, nil
}

// GetEmbeddings returns embeddings for the given input.
func (c *Client) GetEmbeddings(ctx context.Context, input string, opts *common.ChatOptions) ([]float64, error) {
	url := fmt.Sprintf("%s/embeddings", c.ApiBaseURL.Value)
	
	payload := map[string]interface{}{
		"input": input,
		"model": opts.Model,
		// Add other options from opts if supported by LM Studio
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	
	var result struct {
		Data []struct {
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}
	
	return result.Data[0].Embedding, nil
}

// GetName returns the name of the vendor.
func (c *Client) GetName() string {
	return c.Label
}

// IsConfigured checks if the client is configured.
func (c *Client) IsConfigured() bool {
	return c.ApiBaseURL != nil && c.ApiBaseURL.Value != ""
}

// Setup performs any necessary setup for the client.
func (c *Client) Setup() error {
	return c.Configure()
}

// SetupFillEnvFileContent fills the environment file content.
func (c *Client) SetupFillEnvFileContent(buffer *bytes.Buffer) {
	envName := fmt.Sprintf("%s_API_BASE_URL", c.EnvNamePrefix)
	buffer.WriteString(fmt.Sprintf("%s=%s\n", envName, c.ApiBaseURL.Value))
}
