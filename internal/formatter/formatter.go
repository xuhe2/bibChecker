package formatter

import (
	"strings"
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
