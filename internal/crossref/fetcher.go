package crossref

import (
	"context"
	"fmt"
)

// SearchResult represents a simplified search result
type SearchResult struct {
	DOI     string
	Title   string
	Authors string
	Year    string
	Journal string
	URL     string
}

// Fetcher fetches paper data from CrossRef
type Fetcher struct {
	client *Client
}

// NewFetcher creates a new CrossRef Fetcher
func NewFetcher(client *Client) *Fetcher {
	return &Fetcher{client: client}
}

// Search searches for papers by query
func (f *Fetcher) Search(ctx context.Context, query string) ([]SearchResult, error) {
	resp, err := f.client.Search(ctx, query, 20)
	if err != nil {
		return nil, fmt.Errorf("CrossRef search failed: %w", err)
	}

	var results []SearchResult
	for _, work := range resp.Message.Items {
		result := SearchResult{
			DOI:     work.DOI,
			Title:   work.GetTitle(),
			Authors: work.GetAuthorsString(),
			Year:    work.GetYear(),
			Journal: work.GetJournal(),
			URL:     "https://doi.org/" + work.DOI,
		}
		results = append(results, result)
	}

	return results, nil
}

// GetByDOI retrieves paper metadata by DOI
func (f *Fetcher) GetByDOI(ctx context.Context, doi string) (*SearchResult, error) {
	work, err := f.client.GetWork(ctx, doi)
	if err != nil {
		return nil, fmt.Errorf("CrossRef get work failed: %w", err)
	}

	return &SearchResult{
		DOI:     work.DOI,
		Title:   work.GetTitle(),
		Authors: work.GetAuthorsString(),
		Year:    work.GetYear(),
		Journal: work.GetJournal(),
		URL:     "https://doi.org/" + work.DOI,
	}, nil
}

// SearchFirst searches for papers and returns the first result
func (f *Fetcher) SearchFirst(ctx context.Context, query string) (*SearchResult, error) {
	results, err := f.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for query: %s", query)
	}

	return &results[0], nil
}
