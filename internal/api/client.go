package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tehnerd/vape/internal/config"
)

const (
	BaseURL    = "https://atlas.ripe.net/api/v2"
	APITimeout = 30 * time.Second
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: APITimeout,
		},
		baseURL: BaseURL,
		apiKey:  config.GetAPIKey(),
	}
}

func NewClientWithKey(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: APITimeout,
		},
		baseURL: BaseURL,
		apiKey:  apiKey,
	}
}

func (c *Client) doRequest(method, endpoint string, body interface{}, params url.Values) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	urlStr := c.baseURL + endpoint
	if params != nil && len(params) > 0 {
		urlStr += "?" + params.Encode()
	}

	req, err := http.NewRequest(method, urlStr, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Key "+c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, ParseAPIError(resp.StatusCode, respBody)
	}

	return respBody, nil
}

func (c *Client) Get(endpoint string, params url.Values) ([]byte, error) {
	return c.doRequest(http.MethodGet, endpoint, nil, params)
}

func (c *Client) Post(endpoint string, body interface{}) ([]byte, error) {
	return c.doRequest(http.MethodPost, endpoint, body, nil)
}

func (c *Client) Delete(endpoint string) ([]byte, error) {
	return c.doRequest(http.MethodDelete, endpoint, nil, nil)
}
