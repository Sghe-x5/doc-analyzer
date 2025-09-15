package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"local.dev/doc-analyzer/internal/pkg/gateway/handlers"
	pb "local.dev/doc-analyzer/internal/proto/analyzer"
)

// Mock FileAnalysisClient
type MockFileAnalysisClient struct {
	mock.Mock
}

func (m *MockFileAnalysisClient) AnalyzeFile(ctx context.Context, fileID string, generateWordCloud bool) (*pb.AnalyzeFileResponse, error) {
	args := m.Called(ctx, fileID, generateWordCloud)
	return args.Get(0).(*pb.AnalyzeFileResponse), args.Error(1)
}

func (m *MockFileAnalysisClient) GetWordCloud(ctx context.Context, location string) ([]byte, error) {
	args := m.Called(ctx, location)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockFileAnalysisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestAnalyzeFile(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockClient := new(MockFileAnalysisClient)
	handler := handlers.NewAnalysisHandler(mockClient)

	// Create a test server
	router := gin.Default()
	router.POST("/api/v1/analysis", handler.AnalyzeFile)

	// Test case: successful analysis
	t.Run("Successful analysis", func(t *testing.T) {
		// Mock the client response
		mockClient.On("AnalyzeFile", mock.Anything, "file123", true).Return(
			&pb.AnalyzeFileResponse{
				ParagraphCount:    5,
				WordCount:         100,
				CharacterCount:    500,
				IsPlagiarism:      false,
				SimilarFileIds:    []string{},
				WordCloudLocation: "wordclouds/file123.png",
			},
			nil,
		)

		// Create request body
		requestBody := handlers.AnalyzeFileRequest{
			FileID:            "file123",
			GenerateWordCloud: true,
		}
		jsonBody, _ := json.Marshal(requestBody)

		// Create a test request
		req, _ := http.NewRequest("POST", "/api/v1/analysis", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Code)

		// Parse response
		var response handlers.AnalyzeFileResponse
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify response fields
		assert.Equal(t, int32(5), response.ParagraphCount)
		assert.Equal(t, int32(100), response.WordCount)
		assert.Equal(t, int32(500), response.CharacterCount)
		assert.Equal(t, false, response.IsPlagiarism)
		assert.Empty(t, response.SimilarFileIds)
		assert.Equal(t, "wordclouds/file123.png", response.WordCloudLocation)

		mockClient.AssertExpectations(t)
	})

	// Test case: invalid request
	t.Run("Invalid request", func(t *testing.T) {
		// Create an invalid request (missing required field)
		requestBody := map[string]interface{}{
			"generate_word_cloud": true,
		}
		jsonBody, _ := json.Marshal(requestBody)

		// Create a test request
		req, _ := http.NewRequest("POST", "/api/v1/analysis", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockClient.AssertNotCalled(t, "AnalyzeFile")
	})
}

func TestGetWordCloud(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockClient := new(MockFileAnalysisClient)
	handler := handlers.NewAnalysisHandler(mockClient)

	// Create a test server
	router := gin.Default()
	router.GET("/api/v1/wordcloud/:location", handler.GetWordCloud)

	// Test case: successful word cloud retrieval
	t.Run("Successful word cloud retrieval", func(t *testing.T) {
		// Mock image data
		imageData := []byte("fake-image-data")

		// Mock the client response
		mockClient.On("GetWordCloud", mock.Anything, "file123.png").Return(imageData, nil)

		// Create a test request
		req, _ := http.NewRequest("GET", "/api/v1/wordcloud/file123.png", nil)
		resp := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "image/png", resp.Header().Get("Content-Type"))
		assert.Equal(t, imageData, resp.Body.Bytes())

		mockClient.AssertExpectations(t)
	})

	// Test case: missing location parameter
	t.Run("Missing location parameter", func(t *testing.T) {
		// Create a test request with empty location
		req, _ := http.NewRequest("GET", "/api/v1/wordcloud/", nil)
		resp := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, resp.Code) // Gin returns 404 for missing path parameters
		mockClient.AssertNotCalled(t, "GetWordCloud")
	})
}
