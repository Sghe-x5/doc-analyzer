package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"local.dev/doc-analyzer/internal/pkg/analyzer/storage"
)

// MockWordCloudStorage is a mock implementation of the WordCloudStorage interface
type MockWordCloudStorage struct {
	mock.Mock
}

// Ensure MockWordCloudStorage implements WordCloudStorage
var _ storage.WordCloudStorage = (*MockWordCloudStorage)(nil)

// SaveWordCloud mocks the SaveWordCloud method
func (m *MockWordCloudStorage) SaveWordCloud(ctx context.Context, location string, image []byte) error {
	args := m.Called(ctx, location, image)
	return args.Error(0)
}

// GetWordCloud mocks the GetWordCloud method
func (m *MockWordCloudStorage) GetWordCloud(ctx context.Context, location string) ([]byte, error) {
	args := m.Called(ctx, location)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}