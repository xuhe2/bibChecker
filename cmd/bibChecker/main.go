package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"bibChecker/internal/bibtex"
	"bibChecker/internal/config"
	"bibChecker/internal/crossref"
	"bibChecker/internal/formatter"
	"bibChecker/internal/paper"
)

var (
	cfgFile string
	output  string
	noProxy bool
)

var rootCmd = &cobra.Command{
	Use:   "bibChecker",
	Short: "A tool to fetch citations from CrossRef",
}

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search CrossRef and fetch paper metadata",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSearch,
}

var updateCmd = &cobra.Command{
	Use:   "update [bibfile]",
	Short: "Update a BibTeX file with CrossRef metadata",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpdate,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "config file path")
	rootCmd.PersistentFlags().BoolVar(&noProxy, "no-proxy", false, "disable proxy")

	searchCmd.Flags().StringVarP(&output, "output", "o", "", "output file")

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

	client := crossref.NewClient(cfg.Timeout, cfg.CrossRefMailto)
	fetcher := crossref.NewFetcher(client)

	query := strings.Join(args, " ")
	fmt.Printf("Searching: %s\n", query)

	ctx := context.Background()
	results, err := fetcher.Search(ctx, query)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	fmt.Printf("Found %d papers\n\n", len(results))

	var papers []*paper.Paper
	for i, result := range results {
		fmt.Printf("[%d/%d] %s\n", i+1, len(results), result.Title)
		fmt.Printf("  Authors: %s\n", result.Authors)
		fmt.Printf("  Year: %s\n", result.Year)
		fmt.Printf("  Journal: %s\n", result.Journal)
		fmt.Printf("  DOI: %s\n\n", result.DOI)

		p := &paper.Paper{
			Title:   result.Title,
			Authors: strings.Split(result.Authors, " and "),
			Year:    result.Year,
			Journal: result.Journal,
			DOI:     result.DOI,
		}
		papers = append(papers, p)
	}

	// Output results
	outputContent := formatter.JoinPapers(papers)

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

	client := crossref.NewClient(cfg.Timeout, cfg.CrossRefMailto)
	fetcher := crossref.NewFetcher(client)
	bibParser := bibtex.NewParser()

	// Parse input BibTeX file
	inputFile := args[0]
	entries, err := bibParser.ParseFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to parse bib file: %w", err)
	}

	fmt.Printf("Parsed %d entries from %s\n\n", len(entries), inputFile)

	var results []string
	ctx := context.Background()

	for i, entry := range entries {
		fmt.Printf("[%d/%d] %s\n", i+1, len(entries), entry.Title)

		if entry.Title == "" {
			fmt.Printf("  Warning: No title found, keeping original\n")
			results = append(results, entry.RawContent)
			continue
		}

		// Search for the paper on CrossRef
		result, err := fetcher.SearchFirst(ctx, entry.Title)
		if err != nil {
			fmt.Printf("  Warning: search failed: %v\n", err)
			results = append(results, entry.RawContent)
			continue
		}

		fmt.Printf("  Found: %s\n", result.Title)
		fmt.Printf("  DOI: %s\n", result.DOI)

		// Create Paper from search result
		p := &paper.Paper{
			Title:   result.Title,
			Authors: strings.Split(result.Authors, " and "),
			Year:    result.Year,
			Journal: result.Journal,
			DOI:     result.DOI,
		}

		bibtexEntry := p.ToBibTeX()
		if bibtexEntry != "" {
			results = append(results, bibtexEntry)
			fmt.Printf("  Success\n")
		} else {
			fmt.Printf("  Warning: failed to generate BibTeX, keeping original\n")
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
