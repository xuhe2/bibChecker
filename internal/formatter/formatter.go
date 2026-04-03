package formatter

import (
	"strings"

	"bibChecker/internal/paper"
)

// FormatBibTeX formats a BibTeX entry
func FormatBibTeX(entry string) string {
	// Basic formatting - ensure proper line breaks
	lines := strings.Split(entry, "\n")
	var formatted []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			formatted = append(formatted, line)
		}
	}
	return strings.Join(formatted, "\n")
}

// JoinEntries joins multiple BibTeX entries with double newlines
func JoinEntries(entries []string) string {
	var result []string
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry != "" {
			result = append(result, entry)
		}
	}
	return strings.Join(result, "\n\n")
}

// FormatPaper formats a Paper as BibTeX entry
func FormatPaper(p *paper.Paper) string {
	return p.ToBibTeX()
}

// JoinPapers joins multiple Papers as BibTeX entries
func JoinPapers(papers []*paper.Paper) string {
	var entries []string
	for _, p := range papers {
		entries = append(entries, p.ToBibTeX())
	}
	return JoinEntries(entries)
}
