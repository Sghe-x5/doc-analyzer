package analyzer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"local.dev/doc-analyzer/internal/pkg/analyzer/analyzer"
)

func TestTextAnalyzer_AnalyzeText(t *testing.T) {
	// Create a new text analyzer
	textAnalyzer := analyzer.NewTextAnalyzer()

	// Test cases
	testCases := []struct {
		name            string
		content         string
		paragraphCount  int32
		wordCount       int32
		characterCount  int32
	}{
		{
			name:            "Empty text",
			content:         "",
			paragraphCount:  0,
			wordCount:       0,
			characterCount:  0,
		},
		{
			name:            "Single word",
			content:         "Hello",
			paragraphCount:  1,
			wordCount:       1,
			characterCount:  5,
		},
		{
			name:            "Single paragraph",
			content:         "This is a test paragraph with multiple words.",
			paragraphCount:  1,
			wordCount:       8,
			characterCount:  45,
		},
		{
			name:            "Multiple paragraphs",
			content:         "This is the first paragraph.\n\nThis is the second paragraph.\n\nAnd this is the third.",
			paragraphCount:  3,
			wordCount:       15,
			characterCount:  83,
		},
		{
			name:            "Paragraphs with empty lines",
			content:         "Paragraph 1.\n\n\n\nParagraph 2.",
			paragraphCount:  2,
			wordCount:       4,
			characterCount:  28,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			paragraphCount, wordCount, characterCount := textAnalyzer.AnalyzeText(tc.content)

			assert.Equal(t, tc.paragraphCount, paragraphCount, "Paragraph count should match")
			assert.Equal(t, tc.wordCount, wordCount, "Word count should match")
			assert.Equal(t, tc.characterCount, characterCount, "Character count should match")
		})
	}
}

func TestTextAnalyzer_GetNGrams(t *testing.T) {
	// Create a new text analyzer
	textAnalyzer := analyzer.NewTextAnalyzer()

	// Test cases
	testCases := []struct {
		name     string
		content  string
		n        int
		expected []string
	}{
		{
			name:     "Empty text",
			content:  "",
			n:        2,
			expected: []string{},
		},
		{
			name:     "Text shorter than n",
			content:  "Hello world",
			n:        3,
			expected: []string{},
		},
		{
			name:     "Bigrams",
			content:  "This is a test",
			n:        2,
			expected: []string{"This is", "is a", "a test"},
		},
		{
			name:     "Trigrams",
			content:  "This is a test sentence",
			n:        3,
			expected: []string{"This is a", "is a test", "a test sentence"},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := textAnalyzer.GetNGrams(tc.content, tc.n)
			assert.Equal(t, tc.expected, result, "NGrams should match")
		})
	}
}

func TestTextAnalyzer_RemoveExcessWhitespace(t *testing.T) {
	// Create a new text analyzer
	textAnalyzer := analyzer.NewTextAnalyzer()

	// Test cases
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No excess whitespace",
			input:    "This is normal text.",
			expected: "This is normal text.",
		},
		{
			name:     "Multiple spaces",
			input:    "This   has   multiple   spaces.",
			expected: "This has multiple spaces.",
		},
		{
			name:     "Tabs and newlines",
			input:    "This\thas\ttabs\nand\nnewlines.",
			expected: "This has tabs and newlines.",
		},
		{
			name:     "Mixed whitespace",
			input:    "  This  \t has \n mixed \t\n whitespace.  ",
			expected: " This has mixed whitespace. ",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := textAnalyzer.RemoveExcessWhitespace(tc.input)
			assert.Equal(t, tc.expected, result, "Text with removed excess whitespace should match")
		})
	}
}
