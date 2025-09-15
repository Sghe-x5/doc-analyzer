package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"local.dev/doc-analyzer/internal/pkg/storage/storage"
)

// MockFileStorage is a mock implementation of the FileStorage interface
type MockFileStorage struct {
	mock.Mock
}

// Ensure MockFileStorage implements FileStorage
var _ storage.FileStorage = (*MockFileStorage)(nil)

// SaveFile mocks the SaveFile method
func (m *MockFileStorage) SaveFile(ctx context.Context, location string, content []byte) error {
	args := m.Called(ctx, location, content)
	return args.Error(0)
}

// GetFile mocks the GetFile method
func (m *MockFileStorage) GetFile(ctx context.Context, location string) ([]byte, error) {
	args := m.Called(ctx, location)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}