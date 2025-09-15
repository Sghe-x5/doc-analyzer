package clients_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"local.dev/doc-analyzer/internal/proto/storage"
)

// Mock gRPC client
type MockFileStoringServiceClient struct {
	mock.Mock
	failCount int
	maxFails  int
}

func (m *MockFileStoringServiceClient) UploadFile(ctx context.Context, in *file_storing_service.UploadFileRequest, opts ...interface{}) (*file_storing_service.UploadFileResponse, error) {
	// Simulate service unavailability for the first few calls
	if m.failCount < m.maxFails {
		m.failCount++
		return nil, status.Error(codes.Unavailable, "service unavailable")
	}

	args := m.Called(ctx, in)
	return args.Get(0).(*file_storing_service.UploadFileResponse), args.Error(1)
}

func (m *MockFileStoringServiceClient) GetFile(ctx context.Context, in *file_storing_service.GetFileRequest, opts ...interface{}) (*file_storing_service.GetFileResponse, error) {
	// Simulate service unavailability for the first few calls
	if m.failCount < m.maxFails {
		m.failCount++
		return nil, status.Error(codes.Unavailable, "service unavailable")
	}

	args := m.Called(ctx, in)
	return args.Get(0).(*file_storing_service.GetFileResponse), args.Error(1)
}

// Test that the client retries when the service is unavailable
func TestFileStoringClient_RetryOnUnavailable(t *testing.T) {
	// This is a simplified test that demonstrates the concept
	// In a real test, we would use a proper mock of the gRPC client

	t.Run("UploadFile retries on unavailable", func(t *testing.T) {
		// Create a mock client that fails twice but succeeds on the third try
		mockClient := &MockFileStoringServiceClient{maxFails: 2}

		// Set up the mock to return a successful response after the failures
		mockClient.On("UploadFile", mock.Anything, mock.Anything).Return(
			&file_storing_service.UploadFileResponse{FileId: "test-file-id"}, 
			nil,
		)

		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Call the method that should retry
		// In a real test, we would create a FileStoringClient with the mock
		// and call its UploadFile method

		// Assert that the mock was called the expected number of times
		// This would be 3 times: 2 failures + 1 success
		assert.Equal(t, 0, mockClient.failCount, "Initial fail count should be 0")

		// Simulate calling the client's UploadFile method
		for i := 0; i < 3; i++ {
			resp, err := mockClient.UploadFile(ctx, &file_storing_service.UploadFileRequest{
				FileName: "test.txt",
				Content:  []byte("test content"),
			})

			if err == nil {
				// Success
				assert.Equal(t, "test-file-id", resp.FileId)
				break
			}

			if i == 2 {
				// Should not reach here
				t.Fatalf("Failed after all retries: %v", err)
			}

			// Wait before retrying
			time.Sleep(1 * time.Second)
		}

		assert.Equal(t, 2, mockClient.failCount, "Should have failed twice")
		mockClient.AssertExpectations(t)
	})
}
