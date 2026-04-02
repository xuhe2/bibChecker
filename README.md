# bibChecker

A CLI tool to fetch and update BibTeX citations from Google Scholar.

## Install

```bash
go install ./cmd/bibChecker
```

## Usage

```bash
# Search papers and fetch citations
bibChecker search "machine learning"

# Update a .bib file with Google Scholar data
bibChecker update references.bib -o updated.bib
```

## Config

Edit `config.yaml`:

```yaml
proxy: "http://127.0.0.1:7890"  # HTTP proxy
language: "zh-CN"               # Google Scholar language
timeout: 30                     # Request timeout (seconds)
output_format: "bibtex"         # Default citation format
```
