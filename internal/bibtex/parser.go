package bibtex

import (
	"os"
	"regexp"
	"strings"
)

// Entry represents a BibTeX entry
type Entry struct {
	Type       string
	CiteKey    string
	Title      string
	Author     string
	Year       string
	Journal    string
	Volume     string
	Pages      string
	Publisher  string
	Doi        string
	RawContent string
}

// Parser parses BibTeX files
type Parser struct{}

// NewParser creates a new BibTeX parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses a BibTeX content string and returns entries
func (p *Parser) Parse(content string) ([]*Entry, error) {
	var entries []*Entry

	// Match @type{key, ... }
	entryRegex := regexp.MustCompile(`@(\w+)\s*\{\s*([^,]+)\s*,([^@]*)`)
	matches := entryRegex.FindAllStringSubmatch(content, -1)

	for _, m := range matches {
		entry := &Entry{
			Type:       strings.ToLower(m[1]),
			CiteKey:    strings.TrimSpace(m[2]),
			RawContent: m[0],
		}

		// Parse fields from the content
		fields := m[3]

		entry.Title = p.extractField(fields, "title")
		entry.Author = p.extractField(fields, "author")
		entry.Year = p.extractField(fields, "year")
		entry.Journal = p.extractField(fields, "journal")
		entry.Volume = p.extractField(fields, "volume")
		entry.Pages = p.extractField(fields, "pages")
		entry.Publisher = p.extractField(fields, "publisher")
		entry.Doi = p.extractField(fields, "doi")

		entries = append(entries, entry)
	}

	return entries, nil
}

// extractField extracts a field value from BibTeX entry content
func (p *Parser) extractField(content, fieldName string) string {
	// Match field = {value} or field = "value"
	re := regexp.MustCompile(`(?i)` + fieldName + `\s*=\s*[\{"]([^}"]+)[\}"]`)
	m := re.FindStringSubmatch(content)
	if m != nil {
		return strings.TrimSpace(m[1])
	}
	return ""
}

// ParseFile parses a BibTeX file
func (p *Parser) ParseFile(filename string) ([]*Entry, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return p.Parse(string(data))
}
