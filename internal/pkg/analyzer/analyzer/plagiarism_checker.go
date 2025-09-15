package analyzer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// PlagiarismChecker provides methods for checking plagiarism between text documents
// It uses a combination of techniques including:
// 1. Exact matching via hash comparison for efficiency
// 2. N-gram analysis to detect similar text patterns
// 3. Jaccard similarity coefficient to measure text similarity
// 4. Text preprocessing to normalize content before comparison
type PlagiarismChecker struct {
	// Threshold for similarity (0.0 to 1.0)
	// Values closer to 1.0 require higher similarity to be considered plagiarism
	// Default is 0.7 (70% similarity)
	SimilarityThreshold float64

	// Size of n-grams to use for comparison
	// Larger values (4-5) are more specific and reduce false positives
	// Smaller values (2-3) catch more potential matches but may increase false positives
	// Default is 3
	NGramSize int

	// TextAnalyzer instance for word extraction and text processing
	textAnalyzer *TextAnalyzer
}

// NewPlagiarismChecker creates a new PlagiarismChecker instance
func NewPlagiarismChecker() *PlagiarismChecker {
	return &PlagiarismChecker{
		// Default threshold: 30% similarity
		// - Higher values (0.5-0.7) will only detect very similar documents
		// - Lower values (0.2-0.3) will detect documents with minor changes
		// - Values below 0.2 may produce false positives
		SimilarityThreshold: 0.3,

		// Default n-gram size: 3
		// - Larger values (4-5) are more specific and reduce false positives
		// - Smaller values (2-3) catch more potential matches but may increase false positives
		NGramSize:           3,

		textAnalyzer:        NewTextAnalyzer(),
	}
}

// CheckPlagiarism checks if the content is plagiarized from any of the provided contents
// The detection process follows these steps:
// 1. Preprocess the text (remove stop words, normalize whitespace, etc.)
// 2. Generate n-grams from the processed text
// 3. For each comparison text:
//    a. First check for exact matches using hash comparison (fast path)
//    b. If not an exact match, calculate Jaccard similarity between n-gram sets
//    c. If similarity is above the threshold, consider it plagiarism
//
// Parameters:
//   - ctx: Context for the operation
//   - content: The text content to check for plagiarism
//   - otherContents: Map of file IDs to their text contents for comparison
//
// Returns:
//   - bool: True if plagiarism is detected (similarity above threshold)
//   - []string: List of file IDs that are similar to the provided content
func (c *PlagiarismChecker) CheckPlagiarism(ctx context.Context, content string, otherContents map[string]string) (bool, []string) {
	var similarFileIDs []string

	// Preprocess the current content
	processedContent := c.preprocessText(content)

	// Generate n-grams for the current content
	currentNGrams := c.generateNGrams(processedContent, c.NGramSize)

	// Compare with other contents
	for fileID, otherContent := range otherContents {
		// Preprocess the other content
		processedOtherContent := c.preprocessText(otherContent)

		// First, do a quick hash check for exact matches
		if c.calculateHash(processedContent) == c.calculateHash(processedOtherContent) {
			similarFileIDs = append(similarFileIDs, fileID)
			continue
		}

		// If not an exact match, calculate Jaccard similarity
		otherNGrams := c.generateNGrams(processedOtherContent, c.NGramSize)
		similarity := c.calculateJaccardSimilarity(currentNGrams, otherNGrams)

		// Uncomment for debugging
		/*
		if ctx.Value("debug") != nil {
			println("Comparing with", fileID)
			println("Content 1:", processedContent)
			println("Content 2:", processedOtherContent)
			println("Similarity:", similarity)
			println("Threshold:", c.SimilarityThreshold)

			// Print n-grams for debugging
			println("N-grams 1:")
			for ngram := range currentNGrams {
				println("  -", ngram)
			}
			println("N-grams 2:")
			for ngram := range otherNGrams {
				println("  -", ngram)
			}
		}
		*/

		// If similarity is above threshold, consider it plagiarism
		if similarity >= c.SimilarityThreshold {
			similarFileIDs = append(similarFileIDs, fileID)
		}
	}

	return len(similarFileIDs) > 0, similarFileIDs
}

// preprocessText prepares text for comparison by normalizing it
func (c *PlagiarismChecker) preprocessText(text string) string {
	// Get significant words (removes stop words and punctuation)
	significantWords := c.textAnalyzer.GetSignificantWords(text)

	// Join the significant words
	processedText := strings.Join(significantWords, " ")

	// Normalize whitespace
	return c.textAnalyzer.RemoveExcessWhitespace(processedText)
}

// generateNGrams creates a map of n-grams from the text
func (c *PlagiarismChecker) generateNGrams(text string, n int) map[string]int {
	// Use the TextAnalyzer's GetNGrams method to get n-grams
	ngrams := c.textAnalyzer.GetNGrams(text, n)

	// Convert to a frequency map
	ngramFreq := make(map[string]int)
	for _, ngram := range ngrams {
		ngramFreq[ngram]++
	}

	return ngramFreq
}

// calculateJaccardSimilarity computes the Jaccard similarity coefficient between two sets of n-grams
func (c *PlagiarismChecker) calculateJaccardSimilarity(ngrams1, ngrams2 map[string]int) float64 {
	// Create sets from the n-grams
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)

	for ngram := range ngrams1 {
		set1[ngram] = true
	}

	for ngram := range ngrams2 {
		set2[ngram] = true
	}

	// Calculate intersection size
	intersection := 0
	for ngram := range set1 {
		if set2[ngram] {
			intersection++
		}
	}

	// Calculate union size
	union := len(set1) + len(set2) - intersection

	// Avoid division by zero
	if union == 0 {
		return 0
	}

	return float64(intersection) / float64(union)
}

// calculateHash calculates a SHA-256 hash of the content
func (c *PlagiarismChecker) calculateHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}
