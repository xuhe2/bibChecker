package crossref

import (
	"strings"
	"unicode"
)

// NormalizeTitle normalizes a title for comparison
func NormalizeTitle(s string) string {
	s = strings.ToLower(s)
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			return r
		}
		return ' '
	}, s)
	s = strings.Join(strings.Fields(s), " ")
	return strings.TrimSpace(s)
}

// JaroWinkler computes the Jaro-Winkler similarity between two strings (0-1)
func JaroWinkler(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	len1, len2 := len(s1), len(s2)
	if len1 == 0 || len2 == 0 {
		return 0.0
	}

	matchWindow := len1/2 - 1
	if matchWindow < 0 {
		matchWindow = 0
	}

	s1Matches := make([]bool, len1)
	s2Matches := make([]bool, len2)

	matches := 0
	transpositions := 0

	// Find matches
	for i := 0; i < len1; i++ {
		start := i - matchWindow
		if start < 0 {
			start = 0
		}
		end := i + matchWindow + 1
		if end > len2 {
			end = len2
		}

		for j := start; j < end; j++ {
			if s2Matches[j] || s1[i] != s2[j] {
				continue
			}
			s1Matches[i] = true
			s2Matches[j] = true
			matches++
			break
		}
	}

	if matches == 0 {
		return 0.0
	}

	// Count transpositions
	k := 0
	for i := 0; i < len1; i++ {
		if !s1Matches[i] {
			continue
		}
		for !s2Matches[k] {
			k++
		}
		if s1[i] != s2[k] {
			transpositions++
		}
		k++
	}

	// Jaro similarity
	jaro := (float64(matches)/float64(len1) +
		float64(matches)/float64(len2) +
		float64(matches-transpositions)/float64(matches)) / 3.0

	// Winkler modification (prefix bonus)
	prefix := 0
	for i := 0; i < min(4, min(len1, len2)); i++ {
		if s1[i] == s2[i] {
			prefix++
		} else {
			break
		}
	}

	return jaro + float64(prefix)*0.1*(1.0-jaro)
}

// ScoreTitleSimilarity returns a similarity score (0-1) between two titles
func ScoreTitleSimilarity(query, result string) float64 {
	queryNorm := NormalizeTitle(query)
	resultNorm := NormalizeTitle(result)

	// Exact match
	if queryNorm == resultNorm {
		return 1.0
	}

	// Check for exact word overlap
	queryWords := strings.Fields(queryNorm)
	resultWords := strings.Fields(resultNorm)
	commonWords := 0
	for _, qw := range queryWords {
		for _, rw := range resultWords {
			if qw == rw {
				commonWords++
				break
			}
		}
	}

	// Calculate word overlap ratio
	wordOverlap := 0.0
	if len(queryWords) > 0 {
		wordOverlap = float64(commonWords) / float64(len(queryWords))
	}

	// Calculate length ratio
	lengthRatio := float64(len(resultNorm)) / float64(len(queryNorm))
	if lengthRatio > 1 {
		lengthRatio = 1 / lengthRatio // Cap at 1 for the other direction
	}

	// Calculate word coverage - what fraction of result words are in query
	wordCoverage := 0.0
	if len(resultWords) > 0 {
		wordCoverage = float64(commonWords) / float64(len(resultWords))
	}

	// Very high confidence: result is very similar to query with good coverage
	// This handles cases like "U-Net" matching "U-Net: ..."
	if wordCoverage >= 0.9 && wordOverlap >= 0.9 && lengthRatio >= 0.9 {
		return 0.95 + 0.05*JaroWinkler(queryNorm, resultNorm)
	}

	// High confidence: result mostly uses query words and is roughly same length
	if wordCoverage >= 0.8 && wordOverlap >= 0.5 && lengthRatio >= 0.7 {
		return 0.85 + 0.10*JaroWinkler(queryNorm, resultNorm)
	}

	// Good confidence: good coverage, decent overlap, reasonable length
	if wordCoverage >= 0.7 && wordOverlap >= 0.5 && lengthRatio >= 0.6 {
		return 0.80 + 0.08*JaroWinkler(queryNorm, resultNorm)
	}

	// Moderate confidence: decent overlap and coverage
	if wordCoverage >= 0.5 && wordOverlap >= 0.5 && lengthRatio >= 0.5 {
		return 0.70 + 0.05*JaroWinkler(queryNorm, resultNorm)
	}

	// Low confidence: partial match only
	// Apply small JaroWinkler bonus but keep base score low
	lengthPenalty := 0.0
	if lengthRatio < 0.5 {
		lengthPenalty = 0.15
	} else if lengthRatio < 0.7 {
		lengthPenalty = 0.05
	}
	return 0.50 + 0.10*JaroWinkler(queryNorm, resultNorm) - lengthPenalty
}
