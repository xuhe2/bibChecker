package crossref

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Client is the CrossRef API client
type Client struct {
	httpClient *http.Client
	baseURL   string
	mailto    string
}

// NewClient creates a new CrossRef API client
func NewClient(timeout int, mailto string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		baseURL: "https://api.crossref.org/works",
		mailto:  mailto,
	}
}

// Search performs a search query against CrossRef API
func (c *Client) Search(ctx context.Context, query string, rows int) (*SearchResponse, error) {
	if rows <= 0 {
		rows = 20 // default
	}

	params := url.Values{}
	params.Set("query.title", query)
	params.Set("rows", fmt.Sprintf("%d", rows))

	if c.mailto != "" {
		params.Set("mailto", c.mailto)
	}

	// Add random delay to avoid rate limiting
	time.Sleep(time.Duration(100) * time.Millisecond)

	searchURL := c.baseURL + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "bibChecker/1.0 (https://github.com/bibChecker)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetWork retrieves a specific work by DOI
func (c *Client) GetWork(ctx context.Context, doi string) (*Work, error) {
	doi = NormalizeDOI(doi)
	workURL := c.baseURL + "/" + url.PathEscape(doi)

	if c.mailto != "" {
		workURL += "?mailto=" + url.QueryEscape(c.mailto)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", workURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "bibChecker/1.0 (https://github.com/bibChecker)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result struct {
		Message Work `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Message, nil
}
