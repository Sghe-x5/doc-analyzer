package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"local.dev/doc-analyzer/internal/pkg/analyzer/clients"
)

// MockFileStoringClient is a mock implementation of the FileStoringClientInterface
type MockFileStoringClient struct {
	mock.Mock
}

// Ensure MockFileStoringClient implements FileStoringClientInterface
var _ clients.FileStoringClientInterface = (*MockFileStoringClient)(nil)

// GetFile mocks the GetFile method
func (m *MockFileStoringClient) GetFile(ctx context.Context, fileID string) (string, []byte, error) {
	args := m.Called(ctx, fileID)
	return args.String(0), args.Get(1).([]byte), args.Error(2)
}

// Close mocks the Close method
func (m *MockFileStoringClient) Close() error {
	args := m.Called()
	return args.Error(0)
}