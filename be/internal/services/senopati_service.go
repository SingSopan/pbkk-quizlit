package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type SenopatiClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewSenopatiClient creates a new Senopati API client
func NewSenopatiClient() *SenopatiClient {
	baseURL := os.Getenv("SENOPATI_API_BASE_URL")
	if baseURL == "" {
		baseURL = "https://senopati.its.ac.id/senopati-lokal-dev"
	}

	return &SenopatiClient{
		BaseURL: baseURL,
		APIKey:  os.Getenv("SENOPATI_API_KEY"),
		HTTPClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// GenerateRequest represents the request body for /generate endpoint
type GenerateRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Stream      bool    `json:"stream"`
}

// GenerateResponse represents the response from /generate endpoint
type GenerateResponse struct {
	Response string `json:"response"`
	Model    string `json:"model"`
	Done     bool   `json:"done"`
}

// ChatMessage represents a message in the chat conversation
type ChatMessage struct {
	Role    string `json:"role"` // "system", "user", or "assistant"
	Content string `json:"content"`
}

// ChatRequest represents the request body for /chat endpoint
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

// ChatResponse represents the response from /chat endpoint
type ChatResponse struct {
	Message ChatMessage `json:"message"`
	Model   string      `json:"model"`
	Done    bool        `json:"done"`
}

// VisionPDFResponse represents the response from /vision/pdf endpoint
type VisionPDFResponse struct {
	Response string `json:"response"`
	Model    string `json:"model"`
	Done     bool   `json:"done"`
}

// Model represents a model from /models endpoint
type Model struct {
	Name       string `json:"name"`
	ModifiedAt string `json:"modified_at"`
	Size       int64  `json:"size"`
}

// ModelsResponse represents the response from /models endpoint
type ModelsResponse struct {
	Models []string `json:"models"`
}

// GenerateText calls the /generate endpoint
func (c *SenopatiClient) GenerateText(model, prompt string, temperature float64, maxTokens int) (*GenerateResponse, error) {
	reqBody := GenerateRequest{
		Model:       model,
		Prompt:      prompt,
		Temperature: temperature,
		MaxTokens:   maxTokens,
		Stream:      false, // Don't use streaming
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Log request details
	fmt.Printf("[Senopati] POST %s/generate with model=%s, prompt length=%d, temp=%.1f\n",
		c.BaseURL, model, len(prompt), temperature)

	req, err := http.NewRequest("POST", c.BaseURL+"/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("[Senopati] Error response: %s\n", string(body))
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Senopati returns an object like Ollama format
	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Chat calls the /chat endpoint
func (c *SenopatiClient) Chat(model string, messages []ChatMessage, temperature float64, maxTokens int) (*ChatResponse, error) {
	reqBody := ChatRequest{
		Model:       model,
		Messages:    messages,
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// VisionPDF calls the /vision/pdf endpoint with a PDF file
func (c *SenopatiClient) VisionPDF(model, prompt string, pdfData []byte) (*VisionPDFResponse, error) {
	// Create multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add model field
	if err := writer.WriteField("model", model); err != nil {
		return nil, fmt.Errorf("failed to write model field: %w", err)
	}

	// Add prompt field
	if err := writer.WriteField("prompt", prompt); err != nil {
		return nil, fmt.Errorf("failed to write prompt field: %w", err)
	}

	// Add PDF file
	part, err := writer.CreateFormFile("file", "document.pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := part.Write(pdfData); err != nil {
		return nil, fmt.Errorf("failed to write PDF data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", c.BaseURL+"/vision/pdf", &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result VisionPDFResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// ListModels calls the /models endpoint
func (c *SenopatiClient) ListModels() (*ModelsResponse, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/models", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result ModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
