package clients_test

import (
	"context"
	"errors"
	"local.dev/doc-analyzer/internal/pkg/analyzer/clients"
	pb "local.dev/doc-analyzer/internal/proto/storage"
	"local.dev/doc-analyzer/internal/proto/storage/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test implementation of FileStoringClient that uses the mock
type TestFileStoringClient struct {
	clients.FileStoringClientInterface
	mockClient *mocks.MockFileStoringServiceClient
}

func NewTestFileStoringClient(mockClient *mocks.MockFileStoringServiceClient) *TestFileStoringClient {
	return &TestFileStoringClient{
		mockClient: mockClient,
	}
}

func (c *TestFileStoringClient) GetFile(ctx context.Context, fileID string) (string, []byte, error) {
	resp, err := c.mockClient.GetFile(ctx, &pb.GetFileRequest{
		FileId: fileID,
	})

	if err != nil {
		return "", nil, err
	}

	return resp.FileName, resp.Content, nil
}

func (c *TestFileStoringClient) Close() error {
	return nil
}

func TestGetFile(t *testing.T) {
	// Create mock
	mockClient := new(mocks.MockFileStoringServiceClient)

	// Create test client
	client := NewTestFileStoringClient(mockClient)

	// Set up mock expectations
	mockClient.On("GetFile", mock.Anything, &pb.GetFileRequest{
		FileId: "file123",
	}).Return(&pb.GetFileResponse{
		FileName: "test.txt",
		Content:  []byte("test content"),
	}, nil)

	// Call the method
	fileName, content, err := client.GetFile(context.Background(), "file123")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "test.txt", fileName)
	assert.Equal(t, []byte("test content"), content)

	mockClient.AssertExpectations(t)
}

func TestGetFile_Error(t *testing.T) {
	// Create mock
	mockClient := new(mocks.MockFileStoringServiceClient)

	// Create test client
	client := NewTestFileStoringClient(mockClient)

	// Set up mock expectations
	mockClient.On("GetFile", mock.Anything, &pb.GetFileRequest{
		FileId: "file123",
	}).Return(nil, errors.New("connection error"))

	// Call the method
	_, _, err := client.GetFile(context.Background(), "file123")

	// Assert
	assert.Error(t, err)

	mockClient.AssertExpectations(t)
}

func TestClose(t *testing.T) {
	// Create mock
	mockClient := new(mocks.MockFileStoringServiceClient)

	// Create test client
	client := NewTestFileStoringClient(mockClient)

	// Call the method
	err := client.Close()

	// Assert
	assert.NoError(t, err)
}
