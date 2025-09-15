package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// MockFileStoringClient is a mock implementation of the FileStoringClient
type MockFileStoringClient struct {
	mock.Mock
}

// UploadFile mocks the UploadFile method
func (m *MockFileStoringClient) UploadFile(ctx context.Context, fileName string, content []byte) (string, error) {
	args := m.Called(ctx, fileName, content)
	return args.String(0), args.Error(1)
}

// GetFile mocks the GetFile method
func (m *MockFileStoringClient) GetFile(ctx context.Context, fileID string) (string, []byte, error) {
	args := m.Called(ctx, fileID)
	if args.Get(1) == nil {
		return args.String(0), nil, args.Error(2)
	}
	return args.String(0), args.Get(1).([]byte), args.Error(2)
}

// Close mocks the Close method
func (m *MockFileStoringClient) Close() error {
	args := m.Called()
	return args.Error(0)
}