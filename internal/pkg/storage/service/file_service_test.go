package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"local.dev/doc-analyzer/internal/pkg/storage/service"
)

// Mock repository
type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) SaveFile(ctx context.Context, fileID, fileName, hash, location string) error {
	args := m.Called(ctx, fileID, fileName, hash, location)
	return args.Error(0)
}

func (m *MockFileRepository) GetFileByID(ctx context.Context, fileID string) (string, string, error) {
	args := m.Called(ctx, fileID)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockFileRepository) GetFileByHash(ctx context.Context, hash string) (string, error) {
	args := m.Called(ctx, hash)
	return args.String(0), args.Error(1)
}

// Mock storage
type MockFileStorage struct {
	mock.Mock
}

func (m *MockFileStorage) SaveFile(ctx context.Context, location string, content []byte) error {
	args := m.Called(ctx, location, content)
	return args.Error(0)
}

func (m *MockFileStorage) GetFile(ctx context.Context, location string) ([]byte, error) {
	args := m.Called(ctx, location)
	return args.Get(0).([]byte), args.Error(1)
}

func TestFileService_UploadFile(t *testing.T) {
	// Setup
	mockRepo := new(MockFileRepository)
	mockStorage := new(MockFileStorage)
	fileService := service.NewFileService(mockRepo, mockStorage)

	ctx := context.Background()
	fileName := "test.txt"
	content := []byte("test content")

	// Test case: new file upload
	t.Run("New file upload", func(t *testing.T) {
		// Mock repository to return empty fileID (file doesn't exist)
		mockRepo.On("GetFileByHash", ctx, mock.Anything).Return("", nil)
		
		// Mock storage to save file successfully
		mockStorage.On("SaveFile", ctx, mock.Anything, content).Return(nil)
		
		// Mock repository to save file metadata successfully
		mockRepo.On("SaveFile", ctx, mock.Anything, fileName, mock.Anything, mock.Anything).Return(nil)
		
		// Call the method
		fileID, err := fileService.UploadFile(ctx, fileName, content)
		
		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, fileID)
		
		mockRepo.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	// Test case: file already exists
	t.Run("File already exists", func(t *testing.T) {
		// Reset mocks
		mockRepo = new(MockFileRepository)
		mockStorage = new(MockFileStorage)
		fileService = service.NewFileService(mockRepo, mockStorage)

		// Mock repository to return existing fileID
		existingFileID := "existing-file-id"
		mockRepo.On("GetFileByHash", ctx, mock.Anything).Return(existingFileID, nil)
		
		// Call the method
		fileID, err := fileService.UploadFile(ctx, fileName, content)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, existingFileID, fileID)
		
		mockRepo.AssertExpectations(t)
		// Storage should not be called because file already exists
		mockStorage.AssertNotCalled(t, "SaveFile")
	})

	// Test case: error checking file existence
	t.Run("Error checking file existence", func(t *testing.T) {
		// Reset mocks
		mockRepo = new(MockFileRepository)
		mockStorage = new(MockFileStorage)
		fileService = service.NewFileService(mockRepo, mockStorage)

		// Mock repository to return error
		mockRepo.On("GetFileByHash", ctx, mock.Anything).Return("", errors.New("database error"))
		
		// Call the method
		_, err := fileService.UploadFile(ctx, fileName, content)
		
		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check file existence")
		
		mockRepo.AssertExpectations(t)
		// Storage should not be called because of error
		mockStorage.AssertNotCalled(t, "SaveFile")
	})

	// Test case: error saving file content
	t.Run("Error saving file content", func(t *testing.T) {
		// Reset mocks
		mockRepo = new(MockFileRepository)
		mockStorage = new(MockFileStorage)
		fileService = service.NewFileService(mockRepo, mockStorage)

		// Mock repository to return empty fileID (file doesn't exist)
		mockRepo.On("GetFileByHash", ctx, mock.Anything).Return("", nil)
		
		// Mock storage to return error
		mockStorage.On("SaveFile", ctx, mock.Anything, content).Return(errors.New("storage error"))
		
		// Call the method
		_, err := fileService.UploadFile(ctx, fileName, content)
		
		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to save file content")
		
		mockRepo.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
		// Repository should not be called to save metadata because of error
		mockRepo.AssertNotCalled(t, "SaveFile")
	})
}

func TestFileService_GetFile(t *testing.T) {
	// Setup
	mockRepo := new(MockFileRepository)
	mockStorage := new(MockFileStorage)
	fileService := service.NewFileService(mockRepo, mockStorage)

	ctx := context.Background()
	fileID := "file123"

	// Test case: successful file retrieval
	t.Run("Successful file retrieval", func(t *testing.T) {
		// Mock repository to return file metadata
		fileName := "test.txt"
		location := "file123"
		mockRepo.On("GetFileByID", ctx, fileID).Return(fileName, location, nil)
		
		// Mock storage to return file content
		content := []byte("test content")
		mockStorage.On("GetFile", ctx, location).Return(content, nil)
		
		// Call the method
		resultFileName, resultContent, err := fileService.GetFile(ctx, fileID)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, fileName, resultFileName)
		assert.Equal(t, content, resultContent)
		
		mockRepo.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	// Test case: error getting file metadata
	t.Run("Error getting file metadata", func(t *testing.T) {
		// Reset mocks
		mockRepo = new(MockFileRepository)
		mockStorage = new(MockFileStorage)
		fileService = service.NewFileService(mockRepo, mockStorage)

		// Mock repository to return error
		mockRepo.On("GetFileByID", ctx, fileID).Return("", "", errors.New("database error"))
		
		// Call the method
		_, _, err := fileService.GetFile(ctx, fileID)
		
		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get file metadata")
		
		mockRepo.AssertExpectations(t)
		// Storage should not be called because of error
		mockStorage.AssertNotCalled(t, "GetFile")
	})

	// Test case: error getting file content
	t.Run("Error getting file content", func(t *testing.T) {
		// Reset mocks
		mockRepo = new(MockFileRepository)
		mockStorage = new(MockFileStorage)
		fileService = service.NewFileService(mockRepo, mockStorage)

		// Mock repository to return file metadata
		fileName := "test.txt"
		location := "file123"
		mockRepo.On("GetFileByID", ctx, fileID).Return(fileName, location, nil)
		
		// Mock storage to return error
		mockStorage.On("GetFile", ctx, location).Return([]byte{}, errors.New("storage error"))
		
		// Call the method
		_, _, err := fileService.GetFile(ctx, fileID)
		
		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get file content")
		
		mockRepo.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})
}