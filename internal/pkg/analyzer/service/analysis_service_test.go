package service_test

import (
	"context"
	"errors"
	"local.dev/doc-analyzer/internal/pkg/analyzer/clients"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"local.dev/doc-analyzer/internal/pkg/analyzer/analyzer"
	"local.dev/doc-analyzer/internal/pkg/analyzer/service"
)

// Mock repository
type MockAnalysisRepository struct {
	mock.Mock
}

func (m *MockAnalysisRepository) SaveAnalysisResult(ctx context.Context, fileID string, paragraphCount, wordCount, characterCount int32, isPlagiarism bool, wordCloudLocation string) error {
	args := m.Called(ctx, fileID, paragraphCount, wordCount, characterCount, isPlagiarism, wordCloudLocation)
	return args.Error(0)
}

func (m *MockAnalysisRepository) GetAnalysisResult(ctx context.Context, fileID string) (int32, int32, int32, bool, string, error) {
	args := m.Called(ctx, fileID)
	return args.Get(0).(int32), args.Get(1).(int32), args.Get(2).(int32), args.Bool(3), args.String(4), args.Error(5)
}

func (m *MockAnalysisRepository) SaveSimilarFile(ctx context.Context, fileID, similarFileID string) error {
	args := m.Called(ctx, fileID, similarFileID)
	return args.Error(0)
}

func (m *MockAnalysisRepository) GetSimilarFiles(ctx context.Context, fileID string) ([]string, error) {
	args := m.Called(ctx, fileID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAnalysisRepository) GetAllFileIDs(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

// Mock storage
type MockWordCloudStorage struct {
	mock.Mock
}

func (m *MockWordCloudStorage) SaveWordCloud(ctx context.Context, location string, image []byte) error {
	args := m.Called(ctx, location, image)
	return args.Error(0)
}

func (m *MockWordCloudStorage) GetWordCloud(ctx context.Context, location string) ([]byte, error) {
	args := m.Called(ctx, location)
	return args.Get(0).([]byte), args.Error(1)
}

// Mock file storing client
type MockFileStoringClient struct {
	mock.Mock
}

func (m *MockFileStoringClient) GetFile(ctx context.Context, fileID string) (string, []byte, error) {
	args := m.Called(ctx, fileID)
	return args.String(0), args.Get(1).([]byte), args.Error(2)
}

func (m *MockFileStoringClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

// FileStoringClientWrapper wraps the mock to be used as *clients.FileStoringClient
type FileStoringClientWrapper struct {
	client clients.FileStoringClientInterface
}

func TestAnalysisService_AnalyzeFile_ExistingAnalysis(t *testing.T) {
	// Create mocks
	mockRepo := new(MockAnalysisRepository)
	mockStorage := new(MockWordCloudStorage)
	mockFileStoringClient := new(MockFileStoringClient)
	textAnalyzer := analyzer.NewTextAnalyzer()
	plagiarismChecker := analyzer.NewPlagiarismChecker()
	wordCloudGenerator := analyzer.NewWordCloudGenerator("")

	// Create service
	svc := service.NewAnalysisService(
		mockRepo,
		mockStorage,
		mockFileStoringClient,
		textAnalyzer,
		plagiarismChecker,
		wordCloudGenerator,
	)

	// Set up mock expectations for existing analysis
	mockRepo.On("GetAnalysisResult", mock.Anything, "file123").Return(
		int32(5), int32(100), int32(500), false, "wordcloud123.png", nil,
	)

	// Call the method
	paragraphCount, wordCount, characterCount, isPlagiarism, similarFileIDs, wordCloudLocation, err := svc.AnalyzeFile(
		context.Background(), "file123", true,
	)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int32(5), paragraphCount)
	assert.Equal(t, int32(100), wordCount)
	assert.Equal(t, int32(500), characterCount)
	assert.False(t, isPlagiarism)
	assert.Empty(t, similarFileIDs)
	assert.Equal(t, "wordcloud123.png", wordCloudLocation)

	mockRepo.AssertExpectations(t)
	mockStorage.AssertNotCalled(t, "SaveWordCloud")
	mockFileStoringClient.AssertNotCalled(t, "GetFile")
}

func TestAnalysisService_AnalyzeFile_ExistingAnalysisWithPlagiarism(t *testing.T) {
	// Create mocks
	mockRepo := new(MockAnalysisRepository)
	mockStorage := new(MockWordCloudStorage)
	mockFileStoringClient := new(MockFileStoringClient)
	textAnalyzer := analyzer.NewTextAnalyzer()
	plagiarismChecker := analyzer.NewPlagiarismChecker()
	wordCloudGenerator := analyzer.NewWordCloudGenerator("")

	// Create service
	svc := service.NewAnalysisService(
		mockRepo,
		mockStorage,
		mockFileStoringClient,
		textAnalyzer,
		plagiarismChecker,
		wordCloudGenerator,
	)

	// Set up mock expectations for existing analysis with plagiarism
	mockRepo.On("GetAnalysisResult", mock.Anything, "file123").Return(
		int32(5), int32(100), int32(500), true, "wordcloud123.png", nil,
	)
	mockRepo.On("GetSimilarFiles", mock.Anything, "file123").Return(
		[]string{"file456", "file789"}, nil,
	)

	// Call the method
	paragraphCount, wordCount, characterCount, isPlagiarism, similarFileIDs, wordCloudLocation, err := svc.AnalyzeFile(
		context.Background(), "file123", true,
	)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int32(5), paragraphCount)
	assert.Equal(t, int32(100), wordCount)
	assert.Equal(t, int32(500), characterCount)
	assert.True(t, isPlagiarism)
	assert.Equal(t, []string{"file456", "file789"}, similarFileIDs)
	assert.Equal(t, "wordcloud123.png", wordCloudLocation)

	mockRepo.AssertExpectations(t)
	mockStorage.AssertNotCalled(t, "SaveWordCloud")
	mockFileStoringClient.AssertNotCalled(t, "GetFile")
}

func TestAnalysisService_AnalyzeFile_NewAnalysis(t *testing.T) {
	// Create mocks
	mockRepo := new(MockAnalysisRepository)
	mockStorage := new(MockWordCloudStorage)
	mockFileStoringClient := new(MockFileStoringClient)
	textAnalyzer := analyzer.NewTextAnalyzer()
	plagiarismChecker := analyzer.NewPlagiarismChecker()
	wordCloudGenerator := analyzer.NewWordCloudGenerator("")

	// Create service
	svc := service.NewAnalysisService(
		mockRepo,
		mockStorage,
		mockFileStoringClient,
		textAnalyzer,
		plagiarismChecker,
		wordCloudGenerator,
	)

	// Set up mock expectations for new analysis
	mockRepo.On("GetAnalysisResult", mock.Anything, "file123").Return(
		int32(0), int32(0), int32(0), false, "", errors.New("not found"),
	)
	mockFileStoringClient.On("GetFile", mock.Anything, "file123").Return(
		"test.txt", []byte("This is a test file content."), nil,
	)
	mockRepo.On("GetAllFileIDs", mock.Anything).Return(
		[]string{"file456", "file789"}, nil,
	)
	mockFileStoringClient.On("GetFile", mock.Anything, "file456").Return(
		"test456.txt", []byte("This is a different file content."), nil,
	)
	mockFileStoringClient.On("GetFile", mock.Anything, "file789").Return(
		"test789.txt", []byte("This is another file content."), nil,
	)
	mockRepo.On("SaveAnalysisResult", mock.Anything, "file123", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Call the method
	paragraphCount, wordCount, characterCount, isPlagiarism, similarFileIDs, wordCloudLocation, err := svc.AnalyzeFile(
		context.Background(), "file123", false,
	)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int32(1), paragraphCount)
	assert.Equal(t, int32(6), wordCount)
	assert.Equal(t, int32(28), characterCount)
	assert.False(t, isPlagiarism)
	assert.Empty(t, similarFileIDs)
	assert.Empty(t, wordCloudLocation)

	mockRepo.AssertExpectations(t)
	mockStorage.AssertNotCalled(t, "SaveWordCloud")
	mockFileStoringClient.AssertExpectations(t)
}

func TestAnalysisService_AnalyzeFile_NewAnalysisWithWordCloud(t *testing.T) {
	// Create mocks
	mockRepo := new(MockAnalysisRepository)
	mockStorage := new(MockWordCloudStorage)
	mockFileStoringClient := new(MockFileStoringClient)
	textAnalyzer := analyzer.NewTextAnalyzer()
	plagiarismChecker := analyzer.NewPlagiarismChecker()
	wordCloudGenerator := analyzer.NewWordCloudGenerator("")

	// Create service
	svc := service.NewAnalysisService(
		mockRepo,
		mockStorage,
		mockFileStoringClient,
		textAnalyzer,
		plagiarismChecker,
		wordCloudGenerator,
	)

	// Set up mock expectations for new analysis with word cloud
	mockRepo.On("GetAnalysisResult", mock.Anything, "file123").Return(
		int32(0), int32(0), int32(0), false, "", errors.New("not found"),
	)
	mockFileStoringClient.On("GetFile", mock.Anything, "file123").Return(
		"test.txt", []byte("This is a test file content."), nil,
	).Times(2) // Called twice: once for analysis and once for word cloud
	mockRepo.On("GetAllFileIDs", mock.Anything).Return(
		[]string{"file456", "file789"}, nil,
	)
	mockFileStoringClient.On("GetFile", mock.Anything, "file456").Return(
		"test456.txt", []byte("This is a different file content."), nil,
	)
	mockFileStoringClient.On("GetFile", mock.Anything, "file789").Return(
		"test789.txt", []byte("This is another file content."), nil,
	)
	// Mock the word cloud generator to return a test image and location
	mockStorage.On("SaveWordCloud", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(nil)
	mockRepo.On("SaveAnalysisResult", mock.Anything, "file123", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Call the method
	paragraphCount, wordCount, characterCount, isPlagiarism, similarFileIDs, wordCloudLocation, err := svc.AnalyzeFile(
		context.Background(), "file123", true,
	)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int32(1), paragraphCount)
	assert.Equal(t, int32(6), wordCount)
	assert.Equal(t, int32(28), characterCount)
	assert.False(t, isPlagiarism)
	assert.Empty(t, similarFileIDs)
	assert.NotEmpty(t, wordCloudLocation)

	mockRepo.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
	mockFileStoringClient.AssertExpectations(t)
}

func TestAnalysisService_AnalyzeFile_ErrorGettingFile(t *testing.T) {
	// Create mocks
	mockRepo := new(MockAnalysisRepository)
	mockStorage := new(MockWordCloudStorage)
	mockFileStoringClient := new(MockFileStoringClient)
	textAnalyzer := analyzer.NewTextAnalyzer()
	plagiarismChecker := analyzer.NewPlagiarismChecker()
	wordCloudGenerator := analyzer.NewWordCloudGenerator("")

	// Create service
	svc := service.NewAnalysisService(
		mockRepo,
		mockStorage,
		mockFileStoringClient,
		textAnalyzer,
		plagiarismChecker,
		wordCloudGenerator,
	)

	// Set up mock expectations
	mockRepo.On("GetAnalysisResult", mock.Anything, "file123").Return(
		int32(0), int32(0), int32(0), false, "", errors.New("not found"),
	)
	mockFileStoringClient.On("GetFile", mock.Anything, "file123").Return(
		"", []byte(nil), errors.New("file not found"),
	)

	// Call the method
	_, _, _, _, _, _, err := svc.AnalyzeFile(
		context.Background(), "file123", false,
	)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get file content")

	mockRepo.AssertExpectations(t)
	mockFileStoringClient.AssertExpectations(t)
}

func TestAnalysisService_AnalyzeFile_ErrorGettingAllFileIDs(t *testing.T) {
	// Create mocks
	mockRepo := new(MockAnalysisRepository)
	mockStorage := new(MockWordCloudStorage)
	mockFileStoringClient := new(MockFileStoringClient)
	textAnalyzer := analyzer.NewTextAnalyzer()
	plagiarismChecker := analyzer.NewPlagiarismChecker()
	wordCloudGenerator := analyzer.NewWordCloudGenerator("")

	// Create service
	svc := service.NewAnalysisService(
		mockRepo,
		mockStorage,
		mockFileStoringClient,
		textAnalyzer,
		plagiarismChecker,
		wordCloudGenerator,
	)

	// Set up mock expectations
	mockRepo.On("GetAnalysisResult", mock.Anything, "file123").Return(
		int32(0), int32(0), int32(0), false, "", errors.New("not found"),
	)
	mockFileStoringClient.On("GetFile", mock.Anything, "file123").Return(
		"test.txt", []byte("This is a test file content."), nil,
	)
	mockRepo.On("GetAllFileIDs", mock.Anything).Return(
		[]string{}, errors.New("database error"),
	)

	// Call the method
	_, _, _, _, _, _, err := svc.AnalyzeFile(
		context.Background(), "file123", false,
	)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get all file IDs")

	mockRepo.AssertExpectations(t)
	mockFileStoringClient.AssertExpectations(t)
}

func TestAnalysisService_AnalyzeFile_ErrorSavingAnalysisResult(t *testing.T) {
	// Create mocks
	mockRepo := new(MockAnalysisRepository)
	mockStorage := new(MockWordCloudStorage)
	mockFileStoringClient := new(MockFileStoringClient)
	textAnalyzer := analyzer.NewTextAnalyzer()
	plagiarismChecker := analyzer.NewPlagiarismChecker()
	wordCloudGenerator := analyzer.NewWordCloudGenerator("")

	// Create service
	svc := service.NewAnalysisService(
		mockRepo,
		mockStorage,
		mockFileStoringClient,
		textAnalyzer,
		plagiarismChecker,
		wordCloudGenerator,
	)

	// Set up mock expectations
	mockRepo.On("GetAnalysisResult", mock.Anything, "file123").Return(
		int32(0), int32(0), int32(0), false, "", errors.New("not found"),
	)
	mockFileStoringClient.On("GetFile", mock.Anything, "file123").Return(
		"test.txt", []byte("This is a test file content."), nil,
	)
	mockRepo.On("GetAllFileIDs", mock.Anything).Return(
		[]string{"file456", "file789"}, nil,
	)
	mockFileStoringClient.On("GetFile", mock.Anything, "file456").Return(
		"test456.txt", []byte("This is a different file content."), nil,
	)
	mockFileStoringClient.On("GetFile", mock.Anything, "file789").Return(
		"test789.txt", []byte("This is another file content."), nil,
	)
	mockRepo.On("SaveAnalysisResult", mock.Anything, "file123", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("database error"))

	// Call the method
	_, _, _, _, _, _, err := svc.AnalyzeFile(
		context.Background(), "file123", false,
	)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save analysis results")

	mockRepo.AssertExpectations(t)
	mockFileStoringClient.AssertExpectations(t)
}

func TestAnalysisService_AnalyzeFile_ErrorGettingSimilarFiles(t *testing.T) {
	// Create mocks
	mockRepo := new(MockAnalysisRepository)
	mockStorage := new(MockWordCloudStorage)
	mockFileStoringClient := new(MockFileStoringClient)
	textAnalyzer := analyzer.NewTextAnalyzer()
	plagiarismChecker := analyzer.NewPlagiarismChecker()
	wordCloudGenerator := analyzer.NewWordCloudGenerator("")

	// Create service
	svc := service.NewAnalysisService(
		mockRepo,
		mockStorage,
		mockFileStoringClient,
		textAnalyzer,
		plagiarismChecker,
		wordCloudGenerator,
	)

	// Set up mock expectations for existing analysis with plagiarism but error getting similar files
	mockRepo.On("GetAnalysisResult", mock.Anything, "file123").Return(
		int32(5), int32(100), int32(500), true, "wordcloud123.png", nil,
	)
	mockRepo.On("GetSimilarFiles", mock.Anything, "file123").Return(
		[]string{}, errors.New("database error"),
	)

	// Call the method
	_, _, _, _, _, _, err := svc.AnalyzeFile(
		context.Background(), "file123", true,
	)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get similar files")

	mockRepo.AssertExpectations(t)
}

func TestAnalysisService_AnalyzeFile_WithPlagiarism(t *testing.T) {
	// Create mocks
	mockRepo := new(MockAnalysisRepository)
	mockStorage := new(MockWordCloudStorage)
	mockFileStoringClient := new(MockFileStoringClient)
	textAnalyzer := analyzer.NewTextAnalyzer()
	plagiarismChecker := analyzer.NewPlagiarismChecker()
	wordCloudGenerator := analyzer.NewWordCloudGenerator("")

	// Create service
	svc := service.NewAnalysisService(
		mockRepo,
		mockStorage,
		mockFileStoringClient,
		textAnalyzer,
		plagiarismChecker,
		wordCloudGenerator,
	)

	// Set up mock expectations for new analysis with plagiarism
	mockRepo.On("GetAnalysisResult", mock.Anything, "file123").Return(
		int32(0), int32(0), int32(0), false, "", errors.New("not found"),
	)
	mockFileStoringClient.On("GetFile", mock.Anything, "file123").Return(
		"test.txt", []byte("This is a test file content."), nil,
	)
	mockRepo.On("GetAllFileIDs", mock.Anything).Return(
		[]string{"file456"}, nil,
	)
	mockFileStoringClient.On("GetFile", mock.Anything, "file456").Return(
		"test456.txt", []byte("This is a test file content."), nil, // Same content to trigger plagiarism
	)
	mockRepo.On("SaveAnalysisResult", mock.Anything, "file123", mock.Anything, mock.Anything, mock.Anything, true, mock.Anything).Return(nil)
	mockRepo.On("SaveSimilarFile", mock.Anything, "file123", "file456").Return(nil)

	// Call the method
	_, _, _, isPlagiarism, similarFileIDs, _, err := svc.AnalyzeFile(
		context.Background(), "file123", false,
	)

	// Assert
	assert.NoError(t, err)
	assert.True(t, isPlagiarism)
	assert.Contains(t, similarFileIDs, "file456")

	mockRepo.AssertExpectations(t)
	mockFileStoringClient.AssertExpectations(t)
}

func TestAnalysisService_GetWordCloud(t *testing.T) {
	// Create mocks
	mockRepo := new(MockAnalysisRepository)
	mockStorage := new(MockWordCloudStorage)
	mockFileStoringClient := new(MockFileStoringClient)
	textAnalyzer := analyzer.NewTextAnalyzer()
	plagiarismChecker := analyzer.NewPlagiarismChecker()
	wordCloudGenerator := analyzer.NewWordCloudGenerator("")

	// Create service
	svc := service.NewAnalysisService(
		mockRepo,
		mockStorage,
		mockFileStoringClient,
		textAnalyzer,
		plagiarismChecker,
		wordCloudGenerator,
	)

	// Set up mock expectations
	mockStorage.On("GetWordCloud", mock.Anything, "wordcloud123.png").Return(
		[]byte("mock-image-data"), nil,
	)

	// Call the method
	image, err := svc.GetWordCloud(context.Background(), "wordcloud123.png")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []byte("mock-image-data"), image)

	mockStorage.AssertExpectations(t)
}
