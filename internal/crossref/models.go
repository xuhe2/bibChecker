package crossref

import (
	"fmt"
	"strings"
)

// SearchResponse represents the CrossRef API search response
type SearchResponse struct {
	Message struct {
		Items          []Work `json:"items"`
		TotalResults   int    `json:"total-results"`
		ItemsPerPage   int    `json:"items-per-page"`
		Query          struct {
			StartIndex  int    `json:"start-index"`
			ItemsFound  int    `json:"items-found"`
			Time        string `json:"time"`
		} `json:"query"`
	} `json:"message"`
}

// Work represents a work (paper) in CrossRef
type Work struct {
	DOI             string   `json:"DOI"`
	Title           []string `json:"title"`
	Author          []Author `json:"author"`
	Published       Published `json:"published"`
	PublishedPrint  Published `json:"published-print"`
	PublishedOnline Published `json:"published-online"`
	ContainerTitle  []string `json:"container-title"`
	Volume          string   `json:"volume"`
	Issue           string   `json:"issue"`
	Page            string   `json:"page"`
	Publisher       string   `json:"publisher"`
	Type            string   `json:"type"`
	ISSN            []string `json:"ISSN"`
	URL             string   `json:"URL"`
	DOIURL          string   `json:"DOI-url"`
}

// Author represents a paper author
type Author struct {
	Given       string `json:"given"`
	Family      string `json:"family"`
	Sequence    string `json:"sequence"`
	Affiliation []struct {
		Name string `json:"name"`
	} `json:"affiliation"`
}

// Published represents the publication date
type Published struct {
	DateParts [][]int `json:"date-parts"`
}

// GetYear returns the publication year
func (w *Work) GetYear() string {
	// Try published first, then published-print, then published-online
	for _, p := range []Published{w.Published, w.PublishedPrint, w.PublishedOnline} {
		if len(p.DateParts) > 0 && len(p.DateParts[0]) > 0 {
			return fmt.Sprintf("%d", p.DateParts[0][0])
		}
	}
	return ""
}

// GetTitle returns the first title
func (w *Work) GetTitle() string {
	if len(w.Title) > 0 {
		return w.Title[0]
	}
	return ""
}

// GetJournal returns the first container title (journal name)
func (w *Work) GetJournal() string {
	if len(w.ContainerTitle) > 0 {
		return w.ContainerTitle[0]
	}
	return ""
}

// GetAuthorsString returns formatted authors string
func (w *Work) GetAuthorsString() string {
	if len(w.Author) == 0 {
		return ""
	}
	var authors []string
	for _, a := range w.Author {
		if a.Family != "" {
			authors = append(authors, a.Family+", "+a.Given)
		}
	}
	return strings.Join(authors, " and ")
}

// NormalizeDOI normalizes a DOI by removing URL prefix
func NormalizeDOI(doi string) string {
	doi = strings.TrimSpace(doi)
	doi = strings.TrimPrefix(doi, "https://doi.org/")
	doi = strings.TrimPrefix(doi, "http://doi.org/")
	doi = strings.TrimPrefix(doi, "doi.org/")
	// Ensure DOI starts with 10.
	if !strings.HasPrefix(doi, "10.") {
		if strings.HasPrefix(doi, "10.") {
			// Already has 10. prefix
		} else {
			doi = "10." + doi
		}
	}
	return doi
}
