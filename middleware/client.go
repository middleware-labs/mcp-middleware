package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	authHeader string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	return NewClientWithAuth(baseURL, apiKey, "")
}

func NewClientWithAuth(baseURL, apiKey, authorization string) *Client {
	// Normalize baseURL: remove trailing slash if present
	normalizedURL := baseURL
	if len(normalizedURL) > 0 && normalizedURL[len(normalizedURL)-1] == '/' {
		normalizedURL = normalizedURL[:len(normalizedURL)-1]
	}
	return &Client{
		baseURL:    normalizedURL,
		apiKey:     apiKey,
		authHeader: authorization,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) doRequest(ctx context.Context, method, path string, body any, result any) error {
	url := c.baseURL + "/api/v1" + path
	log.Printf("Request: Method=%s Path=%s URL=%s", method, path, url)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		log.Printf("Request Body: %s\n", string(jsonData))
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if c.authHeader != "" {
		req.Header.Set("Authorization", c.authHeader)
	} else if c.apiKey != "" {
		req.Header.Set("ApiKey", c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	// log.Printf("Response Body: %s\n", string(respBody))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
			return fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error)
		}
		// Check if response is HTML (common for error pages)
		bodyStr := string(respBody)
		if len(bodyStr) > 0 && bodyStr[0] == '<' {
			return fmt.Errorf("API error (%d): received HTML response instead of JSON. This usually indicates the endpoint doesn't exist or there's an authentication issue. Response preview: %s", resp.StatusCode, truncateString(bodyStr, 200))
		}
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, truncateString(bodyStr, 500))
	}

	if result != nil && len(respBody) > 0 {
		// Check if response is HTML before trying to unmarshal
		if len(respBody) > 0 && respBody[0] == '<' {
			return fmt.Errorf("received HTML response instead of JSON. This usually indicates the endpoint doesn't exist or there's an authentication issue. Response preview: %s", truncateString(string(respBody), 200))
		}
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
