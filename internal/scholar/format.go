package scholar

// CitationFormat represents the citation format type
type CitationFormat int

const (
	FormatRefWorks CitationFormat = 1
	FormatRIS     CitationFormat = 2
	FormatEndNote  CitationFormat = 3
	FormatBibTeX  CitationFormat = 4
)

func (f CitationFormat) String() string {
	switch f {
	case FormatRefWorks:
		return "1"
	case FormatRIS:
		return "2"
	case FormatEndNote:
		return "3"
	case FormatBibTeX:
		return "4"
	default:
		return "4"
	}
}

// Extension returns the file extension for the format
func (f CitationFormat) Extension() string {
	switch f {
	case FormatRefWorks:
		return ".rfw"
	case FormatRIS:
		return ".ris"
	case FormatEndNote:
		return ".enw"
	case FormatBibTeX:
		return ".bib"
	default:
		return ".bib"
	}
}

// ParseFormat parses a format string to CitationFormat
func ParseFormat(s string) CitationFormat {
	switch s {
	case "bibtex", "bib":
		return FormatBibTeX
	case "ris":
		return FormatRIS
	case "endnote", "enw":
		return FormatEndNote
	case "refworks", "rfw":
		return FormatRefWorks
	default:
		return FormatBibTeX
	}
}
