package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// MockPlagiarismChecker is a mock implementation of the PlagiarismChecker
type MockPlagiarismChecker struct {
	mock.Mock
	SimilarityThreshold float64
	NGramSize           int
}

// CheckPlagiarism mocks the CheckPlagiarism method
func (m *MockPlagiarismChecker) CheckPlagiarism(ctx context.Context, content string, otherContents map[string]string) (bool, []string) {
	args := m.Called(ctx, content, otherContents)
	return args.Bool(0), args.Get(1).([]string)
}

// preprocessText mocks the preprocessText method
func (m *MockPlagiarismChecker) preprocessText(text string) string {
	args := m.Called(text)
	return args.String(0)
}

// generateNGrams mocks the generateNGrams method
func (m *MockPlagiarismChecker) generateNGrams(text string, n int) map[string]int {
	args := m.Called(text, n)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(map[string]int)
}

// calculateJaccardSimilarity mocks the calculateJaccardSimilarity method
func (m *MockPlagiarismChecker) calculateJaccardSimilarity(ngrams1, ngrams2 map[string]int) float64 {
	args := m.Called(ngrams1, ngrams2)
	return args.Get(0).(float64)
}

// calculateHash mocks the calculateHash method
func (m *MockPlagiarismChecker) calculateHash(content string) string {
	args := m.Called(content)
	return args.String(0)
}