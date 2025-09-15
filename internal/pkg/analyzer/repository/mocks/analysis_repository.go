package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"local.dev/doc-analyzer/internal/pkg/analyzer/repository"
)

// MockAnalysisRepository is a mock implementation of the AnalysisRepository interface
type MockAnalysisRepository struct {
	mock.Mock
}

// Ensure MockAnalysisRepository implements AnalysisRepository
var _ repository.AnalysisRepository = (*MockAnalysisRepository)(nil)

// SaveAnalysisResult mocks the SaveAnalysisResult method
func (m *MockAnalysisRepository) SaveAnalysisResult(ctx context.Context, fileID string, paragraphCount, wordCount, characterCount int32, isPlagiarism bool, wordCloudLocation string) error {
	args := m.Called(ctx, fileID, paragraphCount, wordCount, characterCount, isPlagiarism, wordCloudLocation)
	return args.Error(0)
}

// GetAnalysisResult mocks the GetAnalysisResult method
func (m *MockAnalysisRepository) GetAnalysisResult(ctx context.Context, fileID string) (paragraphCount, wordCount, characterCount int32, isPlagiarism bool, wordCloudLocation string, err error) {
	args := m.Called(ctx, fileID)
	return args.Get(0).(int32), args.Get(1).(int32), args.Get(2).(int32), args.Bool(3), args.String(4), args.Error(5)
}

// SaveSimilarFile mocks the SaveSimilarFile method
func (m *MockAnalysisRepository) SaveSimilarFile(ctx context.Context, fileID, similarFileID string) error {
	args := m.Called(ctx, fileID, similarFileID)
	return args.Error(0)
}

// GetSimilarFiles mocks the GetSimilarFiles method
func (m *MockAnalysisRepository) GetSimilarFiles(ctx context.Context, fileID string) ([]string, error) {
	args := m.Called(ctx, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// GetAllFileIDs mocks the GetAllFileIDs method
func (m *MockAnalysisRepository) GetAllFileIDs(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}
