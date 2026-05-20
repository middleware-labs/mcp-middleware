package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL     string
	bearerToken string
	httpClient  *http.Client
}

func NewClient(baseURL, bearerToken string) *Client {
	return &Client{
		baseURL:     strings.TrimRight(baseURL, "/"),
		bearerToken: bearerToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) doRequest(ctx context.Context, method, path string, body any, result any) error {
	if c.baseURL == "" {
		return fmt.Errorf("middleware client: missing tenant base URL")
	}
	if c.bearerToken == "" {
		return fmt.Errorf("middleware client: missing access token")
	}

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

	req.Header.Set("Authorization", "Bearer "+c.bearerToken)
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
			return fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error)
		}
		bodyStr := string(respBody)
		if len(bodyStr) > 0 && bodyStr[0] == '<' {
			return fmt.Errorf("API error (%d): received HTML response instead of JSON. This usually indicates the endpoint doesn't exist or there's an authentication issue. Response preview: %s", resp.StatusCode, truncateString(bodyStr, 200))
		}
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, truncateString(bodyStr, 500))
	}

	if result != nil && len(respBody) > 0 {
		if respBody[0] == '<' {
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
