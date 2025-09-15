package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"local.dev/doc-analyzer/internal/pkg/storage/repository"
)

// MockFileRepository is a mock implementation of the FileRepository interface
type MockFileRepository struct {
	mock.Mock
}

// Ensure MockFileRepository implements FileRepository
var _ repository.FileRepository = (*MockFileRepository)(nil)

// SaveFile mocks the SaveFile method
func (m *MockFileRepository) SaveFile(ctx context.Context, id, name, hash, location string) error {
	args := m.Called(ctx, id, name, hash, location)
	return args.Error(0)
}

// GetFileByID mocks the GetFileByID method
func (m *MockFileRepository) GetFileByID(ctx context.Context, id string) (name string, location string, err error) {
	args := m.Called(ctx, id)
	return args.String(0), args.String(1), args.Error(2)
}

// GetFileByHash mocks the GetFileByHash method
func (m *MockFileRepository) GetFileByHash(ctx context.Context, hash string) (id string, err error) {
	args := m.Called(ctx, hash)
	return args.String(0), args.Error(1)
}