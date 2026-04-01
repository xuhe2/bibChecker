package scholar

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Fetcher handles fetching citation data from Google Scholar
type Fetcher struct {
	client *Client
}

// NewFetcher creates a new Fetcher
func NewFetcher(client *Client) *Fetcher {
	return &Fetcher{client: client}
}

// FetchSearchResults fetches the Google Scholar search results HTML
func (f *Fetcher) FetchSearchResults(query string) (string, error) {
	baseURL := "https://scholar.google.com/scholar"
	params := url.Values{
		"hl":     {f.client.Language()},
		"as_sdt": {"0,5"},
		"q":      {query},
	}

	req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(body), nil
}

// FetchCitePage fetches the cite page for a paper and extracts signature
func (f *Fetcher) FetchCitePage(paperID string) (string, *Signature, error) {
	citeURL := BuildCiteURL(paperID, f.client.Language())

	req, err := http.NewRequest("GET", citeURL, nil)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := f.client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read response: %w", err)
	}

	_, sig := ParseCitePage(string(body))
	return string(body), sig, nil
}

// FetchCitation fetches the citation in the specified format
func (f *Fetcher) FetchCitation(paperID string, format CitationFormat) (string, error) {
	// First get the cite page to extract signature
	_, sig, err := f.FetchCitePage(paperID)
	if err != nil {
		return "", err
	}

	if sig.Scisdr == "" || sig.Scisig == "" {
		return "", fmt.Errorf("failed to extract signature for paper %s", paperID)
	}

	// Build the citation URL
	citationURL := BuildCitationURL(paperID, format, sig, f.client.Language())

	req, err := http.NewRequest("GET", citationURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(body), nil
}

// FetchAllCitations fetches citations in all formats for a paper
func (f *Fetcher) FetchAllCitations(paperID string) (*CitationLinks, error) {
	citeURL := BuildCiteURL(paperID, f.client.Language())

	req, err := http.NewRequest("GET", citeURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	links, _ := ParseCitePage(string(body))
	return links, nil
}
