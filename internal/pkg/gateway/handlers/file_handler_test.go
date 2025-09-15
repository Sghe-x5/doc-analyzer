package handlers

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock FileStoringClient
type MockFileStoringClient struct {
	mock.Mock
}

func (m *MockFileStoringClient) UploadFile(ctx context.Context, fileName string, content []byte) (string, error) {
	args := m.Called(ctx, fileName, content)
	return args.String(0), args.Error(1)
}

func (m *MockFileStoringClient) GetFile(ctx context.Context, fileID string) (string, []byte, error) {
	args := m.Called(ctx, fileID)
	return args.String(0), args.Get(1).([]byte), args.Error(2)
}

func (m *MockFileStoringClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestUploadFile_ValidTxtFile(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockClient := new(MockFileStoringClient)
	handler := NewFileHandler(mockClient)

	// Create a test server
	router := gin.Default()
	router.POST("/api/v1/files", handler.UploadFile)

	// Create a multipart form with a .txt file
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("test content"))
	writer.Close()

	// Mock the client response
	mockClient.On("UploadFile", mock.Anything, "test.txt", mock.Anything).Return("file123", nil)

	// Create a test request
	req, _ := http.NewRequest("POST", "/api/v1/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	mockClient.AssertExpectations(t)
}

func TestUploadFile_InvalidFileExtension(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockClient := new(MockFileStoringClient)
	handler := NewFileHandler(mockClient)

	// Create a test server
	router := gin.Default()
	router.POST("/api/v1/files", handler.UploadFile)

	// Create a multipart form with a non-.txt file
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.pdf")
	part.Write([]byte("test content"))
	writer.Close()

	// Create a test request
	req, _ := http.NewRequest("POST", "/api/v1/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	// The client should not be called because the file extension is invalid
	mockClient.AssertNotCalled(t, "UploadFile")
}

func TestUploadFile_NoFileProvided(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockClient := new(MockFileStoringClient)
	handler := NewFileHandler(mockClient)

	// Create a test server
	router := gin.Default()
	router.POST("/api/v1/files", handler.UploadFile)

	// Create a request with no file
	req, _ := http.NewRequest("POST", "/api/v1/files", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	// The client should not be called because no file was provided
	mockClient.AssertNotCalled(t, "UploadFile")
}

func TestUploadFile_ClientError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockClient := new(MockFileStoringClient)
	handler := NewFileHandler(mockClient)

	// Create a test server
	router := gin.Default()
	router.POST("/api/v1/files", handler.UploadFile)

	// Create a multipart form with a .txt file
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("test content"))
	writer.Close()

	// Mock the client to return an error
	mockClient.On("UploadFile", mock.Anything, "test.txt", mock.Anything).Return("", errors.New("upload error"))

	// Create a test request
	req, _ := http.NewRequest("POST", "/api/v1/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	mockClient.AssertExpectations(t)
}

func TestGetFile_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockClient := new(MockFileStoringClient)
	handler := NewFileHandler(mockClient)

	// Create a test server
	router := gin.Default()
	router.GET("/api/v1/files/:file_id", handler.GetFile)

	// Mock file data
	fileName := "test.txt"
	fileContent := []byte("test content")

	// Mock the client response
	mockClient.On("GetFile", mock.Anything, "file123").Return(fileName, fileContent, nil)

	// Create a test request
	req, _ := http.NewRequest("GET", "/api/v1/files/file123", nil)
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "application/octet-stream", resp.Header().Get("Content-Type"))
	assert.Equal(t, "attachment; filename=test.txt", resp.Header().Get("Content-Disposition"))
	assert.Equal(t, fileContent, resp.Body.Bytes())

	mockClient.AssertExpectations(t)
}

func TestGetFile_MissingFileID(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockClient := new(MockFileStoringClient)
	handler := NewFileHandler(mockClient)

	// Create a test server
	router := gin.Default()
	router.GET("/api/v1/files/:file_id", handler.GetFile)

	// Create a test request with empty file ID
	req, _ := http.NewRequest("GET", "/api/v1/files/", nil)
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.Code) // Gin returns 404 for missing path parameters
	mockClient.AssertNotCalled(t, "GetFile")
}

func TestGetFile_ClientError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockClient := new(MockFileStoringClient)
	handler := NewFileHandler(mockClient)

	// Create a test server
	router := gin.Default()
	router.GET("/api/v1/files/:file_id", handler.GetFile)

	// Mock the client to return an error
	mockClient.On("GetFile", mock.Anything, "file123").Return("", []byte(nil), errors.New("get file error"))

	// Create a test request
	req, _ := http.NewRequest("GET", "/api/v1/files/file123", nil)
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	mockClient.AssertExpectations(t)
}
