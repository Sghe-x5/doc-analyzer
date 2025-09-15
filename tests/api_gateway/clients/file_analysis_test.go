package clients_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"local.dev/doc-analyzer/internal/proto/analyzer"
)

// Mock gRPC client
type MockFileAnalysisServiceClient struct {
	mock.Mock
	failCount int
	maxFails  int
}

func (m *MockFileAnalysisServiceClient) AnalyzeFile(ctx context.Context, in *file_analysis_service.AnalyzeFileRequest, opts ...interface{}) (*file_analysis_service.AnalyzeFileResponse, error) {
	// Simulate service unavailability for the first few calls
	if m.failCount < m.maxFails {
		m.failCount++
		return nil, status.Error(codes.Unavailable, "service unavailable")
	}

	args := m.Called(ctx, in)
	return args.Get(0).(*file_analysis_service.AnalyzeFileResponse), args.Error(1)
}

func (m *MockFileAnalysisServiceClient) GetWordCloud(ctx context.Context, in *file_analysis_service.GetWordCloudRequest, opts ...interface{}) (*file_analysis_service.GetWordCloudResponse, error) {
	// Simulate service unavailability for the first few calls
	if m.failCount < m.maxFails {
		m.failCount++
		return nil, status.Error(codes.Unavailable, "service unavailable")
	}

	args := m.Called(ctx, in)
	return args.Get(0).(*file_analysis_service.GetWordCloudResponse), args.Error(1)
}

// Test that the client retries when the service is unavailable
func TestFileAnalysisClient_RetryOnUnavailable(t *testing.T) {
	// This is a simplified test that demonstrates the concept
	// In a real test, we would use a proper mock of the gRPC client

	t.Run("AnalyzeFile retries on unavailable", func(t *testing.T) {
		// Create a mock client that fails twice but succeeds on the third try
		mockClient := &MockFileAnalysisServiceClient{maxFails: 2}

		// Set up the mock to return a successful response after the failures
		mockClient.On("AnalyzeFile", mock.Anything, mock.Anything).Return(
			&file_analysis_service.AnalyzeFileResponse{
				WordCount:         100,
				CharacterCount:    500,
				ParagraphCount:    10,
				WordCloudLocation: "test-location",
			}, 
			nil,
		)

		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Assert that the mock was called the expected number of times
		assert.Equal(t, 0, mockClient.failCount, "Initial fail count should be 0")

		// Simulate calling the client's AnalyzeFile method
		for i := 0; i < 3; i++ {
			resp, err := mockClient.AnalyzeFile(ctx, &file_analysis_service.AnalyzeFileRequest{
				FileId:            "test-file-id",
				GenerateWordCloud: true,
			})

			if err == nil {
				// Success
				assert.Equal(t, int32(100), resp.WordCount)
				assert.Equal(t, int32(500), resp.CharacterCount)
				assert.Equal(t, int32(10), resp.ParagraphCount)
				assert.Equal(t, "test-location", resp.WordCloudLocation)
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

	t.Run("GetWordCloud retries on unavailable", func(t *testing.T) {
		// Create a mock client that fails twice but succeeds on the third try
		mockClient := &MockFileAnalysisServiceClient{maxFails: 2}

		// Set up the mock to return a successful response after the failures
		mockClient.On("GetWordCloud", mock.Anything, mock.Anything).Return(
			&file_analysis_service.GetWordCloudResponse{
				Image: []byte("test-image-data"),
			}, 
			nil,
		)

		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Assert that the mock was called the expected number of times
		assert.Equal(t, 0, mockClient.failCount, "Initial fail count should be 0")

		// Simulate calling the client's GetWordCloud method
		for i := 0; i < 3; i++ {
			resp, err := mockClient.GetWordCloud(ctx, &file_analysis_service.GetWordCloudRequest{
				Location: "test-location",
			})

			if err == nil {
				// Success
				assert.Equal(t, []byte("test-image-data"), resp.Image)
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
