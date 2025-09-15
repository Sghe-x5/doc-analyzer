package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	pb "local.dev/doc-analyzer/internal/proto/analyzer"
)

// MockFileAnalysisClient is a mock implementation of the FileAnalysisClient
type MockFileAnalysisClient struct {
	mock.Mock
}

// AnalyzeFile mocks the AnalyzeFile method
func (m *MockFileAnalysisClient) AnalyzeFile(ctx context.Context, fileID string, generateWordCloud bool) (*pb.AnalyzeFileResponse, error) {
	args := m.Called(ctx, fileID, generateWordCloud)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.AnalyzeFileResponse), args.Error(1)
}

// GetWordCloud mocks the GetWordCloud method
func (m *MockFileAnalysisClient) GetWordCloud(ctx context.Context, location string) ([]byte, error) {
	args := m.Called(ctx, location)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

// Close mocks the Close method
func (m *MockFileAnalysisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}