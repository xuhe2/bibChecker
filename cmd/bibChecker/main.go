package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"bibChecker/internal/bibtex"
	"bibChecker/internal/config"
	"bibChecker/internal/formatter"
	"bibChecker/internal/scholar"
)

var (
	cfgFile  string
	output   string
	format   string
	noProxy  bool
)

var rootCmd = &cobra.Command{
	Use:   "bibChecker",
	Short: "A tool to fetch citations from Google Scholar",
}

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search Google Scholar and fetch citations",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSearch,
}

var updateCmd = &cobra.Command{
	Use:   "update [bibfile]",
	Short: "Update a BibTeX file with Google Scholar data",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpdate,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "config file path")
	rootCmd.PersistentFlags().BoolVar(&noProxy, "no-proxy", false, "disable proxy")

	searchCmd.Flags().StringVarP(&output, "output", "o", "", "output file")
	searchCmd.Flags().StringVarP(&format, "format", "f", "bibtex", "citation format (bibtex/ris/endnote/refworks)")

	updateCmd.Flags().StringVarP(&output, "output", "o", "", "output file (default: input.updated.bib)")

	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(updateCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func loadConfig() (*config.Config, error) {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return nil, err
	}

	if noProxy {
		cfg.Proxy = ""
	}

	return cfg, nil
}

func runSearch(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := scholar.NewClient(cfg.Proxy, cfg.Language, cfg.Timeout)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	fetcher := scholar.NewFetcher(client)

	query := strings.Join(args, " ")
	fmt.Printf("Searching: %s\n", query)

	// Fetch search results
	html, err := fetcher.FetchSearchResults(query)
	if err != nil {
		return fmt.Errorf("failed to fetch search results: %w", err)
	}

	// Parse paper IDs
	paperIDs := scholar.ParseSearchResults(html)
	fmt.Printf("Found %d papers\n\n", len(paperIDs))

	citationFormat := scholar.ParseFormat(format)

	var results []string
	for i, paperID := range paperIDs {
		fmt.Printf("[%d/%d] Fetching citation for %s...\n", i+1, len(paperIDs), paperID)

		citation, err := fetcher.FetchCitation(paperID, citationFormat)
		if err != nil {
			fmt.Printf("  Warning: %v\n", err)
			continue
		}

		if citation != "" {
			results = append(results, citation)
			fmt.Printf("  Success (%d bytes)\n", len(citation))
		}
	}

	// Output results
	outputContent := formatter.JoinEntries(results)

	if output != "" {
		if err := os.WriteFile(output, []byte(outputContent), 0644); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		fmt.Printf("\nOutput written to: %s\n", output)
	} else {
		fmt.Println("\n" + outputContent)
	}

	return nil
}

func runUpdate(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := scholar.NewClient(cfg.Proxy, cfg.Language, cfg.Timeout)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	fetcher := scholar.NewFetcher(client)
	bibParser := bibtex.NewParser()

	// Parse input BibTeX file
	inputFile := args[0]
	entries, err := bibParser.ParseFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to parse bib file: %w", err)
	}

	fmt.Printf("Parsed %d entries from %s\n\n", len(entries), inputFile)

	var results []string
	for i, entry := range entries {
		fmt.Printf("[%d/%d] Processing: %s\n", i+1, len(entries), entry.Title)

		if entry.Title == "" {
			fmt.Printf("  Warning: No title found, keeping original\n")
			results = append(results, entry.RawContent)
			continue
		}

		// Search for the paper
		html, err := fetcher.FetchSearchResults(entry.Title)
		if err != nil {
			fmt.Printf("  Warning: search failed: %v\n", err)
			results = append(results, entry.RawContent)
			continue
		}

		paperIDs := scholar.ParseSearchResults(html)
		if len(paperIDs) == 0 {
			fmt.Printf("  Warning: no results found, keeping original\n")
			results = append(results, entry.RawContent)
			continue
		}

		// Use the first result
		paperID := paperIDs[0]
		fmt.Printf("  Found paper ID: %s\n", paperID)

		citation, err := fetcher.FetchCitation(paperID, scholar.FormatBibTeX)
		if err != nil {
			fmt.Printf("  Warning: fetch failed: %v\n", err)
			results = append(results, entry.RawContent)
			continue
		}

		if citation != "" && strings.Contains(citation, "@") {
			results = append(results, citation)
			fmt.Printf("  Success\n")
		} else {
			fmt.Printf("  Warning: invalid citation, keeping original\n")
			results = append(results, entry.RawContent)
		}
	}

	// Determine output file
	if output == "" {
		output = strings.TrimSuffix(inputFile, ".bib") + ".updated.bib"
	}

	// Write output
	outputContent := formatter.JoinEntries(results)
	if err := os.WriteFile(output, []byte(outputContent), 0644); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	fmt.Printf("\nUpdated BibTeX written to: %s\n", output)
	return nil
}
