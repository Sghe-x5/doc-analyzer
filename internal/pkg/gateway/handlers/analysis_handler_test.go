package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

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

func TestAnalyzeFile_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockClient := new(MockFileAnalysisClient)
	handler := NewAnalysisHandler(mockClient)

	// Create a test server
	router := gin.Default()
	router.POST("/api/v1/analysis", handler.AnalyzeFile)

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
	requestBody := AnalyzeFileRequest{
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
	var response AnalyzeFileResponse
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
}

func TestAnalyzeFile_InvalidRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockClient := new(MockFileAnalysisClient)
	handler := NewAnalysisHandler(mockClient)

	// Create a test server
	router := gin.Default()
	router.POST("/api/v1/analysis", handler.AnalyzeFile)

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
}

func TestAnalyzeFile_ClientError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockClient := new(MockFileAnalysisClient)
	handler := NewAnalysisHandler(mockClient)

	// Create a test server
	router := gin.Default()
	router.POST("/api/v1/analysis", handler.AnalyzeFile)

	// Mock the client to return an error
	mockClient.On("AnalyzeFile", mock.Anything, "file123", true).Return(
		&pb.AnalyzeFileResponse{},
		errors.New("analysis error"),
	)

	// Create request body
	requestBody := AnalyzeFileRequest{
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
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	mockClient.AssertExpectations(t)
}

func TestGetWordCloud_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockClient := new(MockFileAnalysisClient)
	handler := NewAnalysisHandler(mockClient)

	// Create a test server
	router := gin.Default()
	router.GET("/api/v1/wordcloud/:location", handler.GetWordCloud)

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
}

func TestGetWordCloud_MissingLocation(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockClient := new(MockFileAnalysisClient)
	handler := NewAnalysisHandler(mockClient)

	// Create a test server
	router := gin.Default()
	router.GET("/api/v1/wordcloud/:location", handler.GetWordCloud)

	// Create a test request with empty location
	req, _ := http.NewRequest("GET", "/api/v1/wordcloud/", nil)
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.Code) // Gin returns 404 for missing path parameters
	mockClient.AssertNotCalled(t, "GetWordCloud")
}

func TestGetWordCloud_ClientError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockClient := new(MockFileAnalysisClient)
	handler := NewAnalysisHandler(mockClient)

	// Create a test server
	router := gin.Default()
	router.GET("/api/v1/wordcloud/:location", handler.GetWordCloud)

	// Mock the client to return an error
	mockClient.On("GetWordCloud", mock.Anything, "file123.png").Return([]byte(nil), errors.New("get wordcloud error"))

	// Create a test request
	req, _ := http.NewRequest("GET", "/api/v1/wordcloud/file123.png", nil)
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	mockClient.AssertExpectations(t)
}
