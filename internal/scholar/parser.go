package scholar

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Paper represents a paper from Google Scholar search results
type Paper struct {
	ID       string // data-cid attribute
	Title    string
	URL      string
	Year     string
	Authors  string
	Journal  string
	CitedBy  int
}

// CitationLinks contains the citation download URLs for different formats
type CitationLinks struct {
	BibTeX    string
	EndNote   string
	RIS       string
	RefWorks  string
}

// Signature contains scisdr and scisig parameters
type Signature struct {
	Scisdr string
	Scisig string
}

// ParseSearchResults extracts paper IDs from search result HTML using goquery
func ParseSearchResults(html string) []string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil
	}

	var ids []string
	seen := make(map[string]bool)

	// First find #gs_bdy, then find descendants with data-cid attribute
	doc.Find("#gs_bdy [data-cid]").Each(func(i int, s *goquery.Selection) {
		cid, exists := s.Attr("data-cid")
		if exists && cid != "" && !seen[cid] {
			seen[cid] = true
			ids = append(ids, cid)
		}
	})

	return ids
}

// ParseCitePage extracts citation links and signature from cite page HTML using goquery
func ParseCitePage(html string) (*CitationLinks, *Signature) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, nil
	}

	links := &CitationLinks{}
	sig := &Signature{}

	// Extract citation links by finding <a> tags containing specific text
	doc.Find("a.gs_citi").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		text := strings.TrimSpace(s.Text())

		if strings.Contains(href, "scholar.bib") && text == "BibTeX" {
			links.BibTeX = href
		} else if strings.Contains(href, "scholar.enw") && text == "EndNote" {
			links.EndNote = href
		} else if strings.Contains(href, "scholar.ris") && text == "RefMan" {
			links.RIS = href
		} else if strings.Contains(href, "scholar.rfw") && text == "RefWorks" {
			links.RefWorks = href
		}
	})

	// Extract scisdr and scisig from the first citation link
	if links.BibTeX != "" {
		parsedURL, err := url.Parse(links.BibTeX)
		if err == nil {
			query := parsedURL.Query()
			sig.Scisdr = query.Get("scisdr")
			sig.Scisig = query.Get("scisig")
		}
	}

	return links, sig
}

// BuildCiteURL constructs the cite page URL for a paper
func BuildCiteURL(paperID, language string) string {
	return fmt.Sprintf("https://scholar.google.com/scholar?q=info:%s:scholar.google.com/&output=cite&scirp=0&hl=%s", paperID, language)
}

// BuildCitationURL constructs a citation download URL with signature
func BuildCitationURL(paperID string, format CitationFormat, sig *Signature, language string) string {
	base := "https://scholar.googleusercontent.com/scholar."

	switch format {
	case FormatBibTeX:
		base += "bib"
	case FormatEndNote:
		base += "enw"
	case FormatRIS:
		base += "ris"
	case FormatRefWorks:
		base += "rfw"
	}

	return base + fmt.Sprintf("?q=info:%s:scholar.google.com/&output=citation&scisdr=%s&scisig=%s&scisf=%s&ct=citation&cd=-1&hl=%s",
		paperID, sig.Scisdr, sig.Scisig, format.String(), language)
}
