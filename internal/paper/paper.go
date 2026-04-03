package paper

import (
	"fmt"
	"regexp"
	"strings"
)

// Paper represents a research paper with metadata
type Paper struct {
	DOI       string
	Title     string
	Authors   []string
	Year      string
	Journal   string
	Volume    string
	Issue     string
	Pages     string
	Publisher string
}

// FormatAuthors formats authors as "Last, First and Last2, First2"
func (p *Paper) FormatAuthors() string {
	if len(p.Authors) == 0 {
		return ""
	}
	var formatted []string
	for _, a := range p.Authors {
		formatted = append(formatted, a)
	}
	return strings.Join(formatted, " and ")
}

// GenerateCiteKey generates a citation key from title and year
func (p *Paper) GenerateCiteKey() string {
	// Extract first meaningful word from title
	title := strings.ToLower(p.Title)
	title = regexp.MustCompile(`[^a-z]+`).ReplaceAllString(title, " ")
	title = strings.Trim(title, " ")
	words := strings.Split(title, " ")
	key := ""
	if len(words) > 0 && words[0] != "" {
		key = words[0]
	}
	if len(words) > 1 && words[1] != "" {
		key += words[1]
	}

	// Add year
	if p.Year != "" {
		key += p.Year
	}
	return key
}

// ToBibTeX generates a BibTeX entry string
func (p *Paper) ToBibTeX() string {
	citeKey := p.GenerateCiteKey()
	entryType := "article"
	if p.Journal == "" {
		entryType = "misc"
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("@%s{%s,", entryType, citeKey))

	if p.Title != "" {
		lines = append(lines, fmt.Sprintf("  title = {%s},", p.Title))
	}
	if p.FormatAuthors() != "" {
		lines = append(lines, fmt.Sprintf("  author = {%s},", p.FormatAuthors()))
	}
	if p.Year != "" {
		lines = append(lines, fmt.Sprintf("  year = {%s},", p.Year))
	}
	if p.Journal != "" {
		lines = append(lines, fmt.Sprintf("  journal = {%s},", p.Journal))
	}
	if p.Volume != "" {
		lines = append(lines, fmt.Sprintf("  volume = {%s},", p.Volume))
	}
	if p.Issue != "" {
		lines = append(lines, fmt.Sprintf("  number = {%s},", p.Issue))
	}
	if p.Pages != "" {
		lines = append(lines, fmt.Sprintf("  pages = {%s},", p.Pages))
	}
	if p.DOI != "" {
		lines = append(lines, fmt.Sprintf("  doi = {%s},", p.DOI))
	}
	if p.Publisher != "" {
		lines = append(lines, fmt.Sprintf("  publisher = {%s},", p.Publisher))
	}

	// Remove trailing comma from last line
	if len(lines) > 1 {
		lastIdx := len(lines) - 1
		lines[lastIdx] = strings.TrimRight(lines[lastIdx], ",")
	}

	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}
