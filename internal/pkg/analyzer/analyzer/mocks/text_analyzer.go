package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockTextAnalyzer is a mock implementation of the TextAnalyzer
type MockTextAnalyzer struct {
	mock.Mock
	StopWords map[string]bool
}

// AnalyzeText mocks the AnalyzeText method
func (m *MockTextAnalyzer) AnalyzeText(content string) (paragraphCount, wordCount, characterCount int32) {
	args := m.Called(content)
	return args.Get(0).(int32), args.Get(1).(int32), args.Get(2).(int32)
}

// GetWords mocks the GetWords method
func (m *MockTextAnalyzer) GetWords(content string) []string {
	args := m.Called(content)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]string)
}

// GetSignificantWords mocks the GetSignificantWords method
func (m *MockTextAnalyzer) GetSignificantWords(content string) []string {
	args := m.Called(content)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]string)
}

// GetNGrams mocks the GetNGrams method
func (m *MockTextAnalyzer) GetNGrams(content string, n int) []string {
	args := m.Called(content, n)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]string)
}

// RemoveExcessWhitespace mocks the RemoveExcessWhitespace method
func (m *MockTextAnalyzer) RemoveExcessWhitespace(text string) string {
	args := m.Called(text)
	return args.String(0)
}