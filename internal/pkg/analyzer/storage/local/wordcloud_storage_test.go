package local_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"local.dev/doc-analyzer/internal/pkg/analyzer/storage/local"
)

func TestNewLocalStorage(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "wordcloud_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test successful creation
	storage, err := local.NewLocalStorage(tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, storage)

	// Test creation with non-existent directory (should create it)
	nonExistentDir := filepath.Join(tempDir, "non_existent")
	storage, err = local.NewLocalStorage(nonExistentDir)
	assert.NoError(t, err)
	assert.NotNil(t, storage)

	// Verify the directory was created
	_, err = os.Stat(nonExistentDir)
	assert.NoError(t, err)

	// Test creation with empty path (should use current directory)
	t.Run("Empty path", func(t *testing.T) {
		storage, err := local.NewLocalStorage("")
		assert.NoError(t, err)
		assert.NotNil(t, storage)
	})

	// Test creation with relative path
	t.Run("Relative path", func(t *testing.T) {
		relativeDir := "./test_relative_dir"
		defer os.RemoveAll(relativeDir)

		storage, err := local.NewLocalStorage(relativeDir)
		assert.NoError(t, err)
		assert.NotNil(t, storage)

		// Verify the directory was created
		_, err = os.Stat(relativeDir)
		assert.NoError(t, err)
	})

	// Test error when creating directory fails
	t.Run("MkdirAll error", func(t *testing.T) {
		// Create a file with the same name as the directory we'll try to create
		// This will cause MkdirAll to fail because a file exists with that name
		filePath := filepath.Join(tempDir, "file-as-dir")
		err := os.WriteFile(filePath, []byte("test"), 0644)
		require.NoError(t, err)

		// Try to create a storage with a path that includes the file as a directory
		invalidPath := filepath.Join(filePath, "subdir")
		storage, err := local.NewLocalStorage(invalidPath)

		// Should fail because we can't create a directory with the same name as a file
		assert.Error(t, err)
		assert.Nil(t, storage)
		assert.Contains(t, err.Error(), "failed to create storage directory")
	})
}

func TestSaveWordCloud(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "wordcloud_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a storage instance
	storage, err := local.NewLocalStorage(tempDir)
	require.NoError(t, err)

	// Test saving a word cloud
	testImage := []byte("test image data")
	location := "test.png"
	err = storage.SaveWordCloud(context.Background(), location, testImage)
	assert.NoError(t, err)

	// Verify the file was created
	fullPath := filepath.Join(tempDir, location)
	savedImage, err := os.ReadFile(fullPath)
	assert.NoError(t, err)
	assert.Equal(t, testImage, savedImage)

	// Test saving to a subdirectory
	subDirLocation := filepath.Join("subdir", "test.png")
	err = storage.SaveWordCloud(context.Background(), subDirLocation, testImage)
	assert.NoError(t, err)

	// Verify the file was created in the subdirectory
	fullSubDirPath := filepath.Join(tempDir, subDirLocation)
	savedSubDirImage, err := os.ReadFile(fullSubDirPath)
	assert.NoError(t, err)
	assert.Equal(t, testImage, savedSubDirImage)

	// Test saving with empty location
	t.Run("Empty location", func(t *testing.T) {
		err = storage.SaveWordCloud(context.Background(), "", testImage)
		assert.NoError(t, err)

		// Verify the file was created with the default filename
		defaultPath := filepath.Join(tempDir, "default.png")
		savedEmptyLocImage, err := os.ReadFile(defaultPath)
		assert.NoError(t, err)
		assert.Equal(t, testImage, savedEmptyLocImage)
	})

	// Test saving with empty content
	t.Run("Empty content", func(t *testing.T) {
		emptyImage := []byte{}
		emptyContentLoc := "empty.png"
		err = storage.SaveWordCloud(context.Background(), emptyContentLoc, emptyImage)
		assert.NoError(t, err)

		// Verify the empty file was created
		emptyContentPath := filepath.Join(tempDir, emptyContentLoc)
		savedEmptyImage, err := os.ReadFile(emptyContentPath)
		assert.NoError(t, err)
		assert.Empty(t, savedEmptyImage)
	})

	// Test overwriting existing file
	t.Run("Overwrite existing file", func(t *testing.T) {
		// First save
		firstImage := []byte("first image data")
		overwriteLoc := "overwrite.png"
		err = storage.SaveWordCloud(context.Background(), overwriteLoc, firstImage)
		assert.NoError(t, err)

		// Verify first save
		overwritePath := filepath.Join(tempDir, overwriteLoc)
		savedFirstImage, err := os.ReadFile(overwritePath)
		assert.NoError(t, err)
		assert.Equal(t, firstImage, savedFirstImage)

		// Second save (overwrite)
		secondImage := []byte("second image data")
		err = storage.SaveWordCloud(context.Background(), overwriteLoc, secondImage)
		assert.NoError(t, err)

		// Verify overwrite
		savedSecondImage, err := os.ReadFile(overwritePath)
		assert.NoError(t, err)
		assert.Equal(t, secondImage, savedSecondImage)
	})

	// Test error when creating directory fails
	t.Run("MkdirAll error", func(t *testing.T) {
		// Create a file with the same name as the directory we'll try to create
		filePath := filepath.Join(tempDir, "file-as-dir-save")
		err := os.WriteFile(filePath, []byte("test"), 0644)
		require.NoError(t, err)

		// Try to save a file in a subdirectory of the file (which will fail)
		badLocation := filepath.Join("file-as-dir-save", "test.png")
		err = storage.SaveWordCloud(context.Background(), badLocation, testImage)

		// Should fail because we can't create a directory with the same name as a file
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create directory")
	})

	// Test error when writing file fails
	t.Run("WriteFile error", func(t *testing.T) {
		// Create a directory with the same name as the file we'll try to create
		dirPath := filepath.Join(tempDir, "dir-as-file")
		err := os.MkdirAll(dirPath, 0755)
		require.NoError(t, err)

		// Try to save a file with the same name as the directory (which will fail)
		err = storage.SaveWordCloud(context.Background(), "dir-as-file", testImage)

		// Should fail because we can't create a file with the same name as a directory
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to write file")
	})
}

func TestGetWordCloud(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "wordcloud_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a storage instance
	storage, err := local.NewLocalStorage(tempDir)
	require.NoError(t, err)

	// Create a test image file
	testImage := []byte("test image data")
	location := "test.png"
	fullPath := filepath.Join(tempDir, location)
	err = os.WriteFile(fullPath, testImage, 0644)
	require.NoError(t, err)

	// Test retrieving the word cloud
	retrievedImage, err := storage.GetWordCloud(context.Background(), location)
	assert.NoError(t, err)
	assert.Equal(t, testImage, retrievedImage)

	// Test retrieving a non-existent word cloud
	_, err = storage.GetWordCloud(context.Background(), "non_existent.png")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Test retrieving with empty location
	t.Run("Empty location", func(t *testing.T) {
		// Create a file with default filename
		defaultPath := filepath.Join(tempDir, "default.png")
		defaultImage := []byte("default image")
		err = os.WriteFile(defaultPath, defaultImage, 0644)
		require.NoError(t, err)

		// Retrieve the file with empty location (should use default filename)
		retrievedDefaultImage, err := storage.GetWordCloud(context.Background(), "")
		assert.NoError(t, err)
		assert.Equal(t, defaultImage, retrievedDefaultImage)
	})

	// Test retrieving an empty file
	t.Run("Empty file", func(t *testing.T) {
		// Create an empty file
		emptyFilePath := filepath.Join(tempDir, "empty.png")
		err = os.WriteFile(emptyFilePath, []byte{}, 0644)
		require.NoError(t, err)

		// Retrieve the empty file
		retrievedEmptyFile, err := storage.GetWordCloud(context.Background(), "empty.png")
		assert.NoError(t, err)
		assert.Empty(t, retrievedEmptyFile)
	})

	// Test error when reading file fails (not due to file not existing)
	t.Run("ReadFile error", func(t *testing.T) {
		// Create a directory with the same name as the file we'll try to read
		dirPath := filepath.Join(tempDir, "dir-as-file-read")
		err := os.MkdirAll(dirPath, 0755)
		require.NoError(t, err)

		// Try to read a "file" with the same name as the directory (which will fail)
		_, err = storage.GetWordCloud(context.Background(), "dir-as-file-read")

		// Should fail because we can't read a directory as a file
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read file")
	})
}

func TestSaveAndGetWordCloud(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "wordcloud_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a storage instance
	storage, err := local.NewLocalStorage(tempDir)
	require.NoError(t, err)

	// Test saving and then retrieving a word cloud
	testImage := []byte("test image data")
	location := "test.png"

	// Save the word cloud
	err = storage.SaveWordCloud(context.Background(), location, testImage)
	assert.NoError(t, err)

	// Retrieve the word cloud
	retrievedImage, err := storage.GetWordCloud(context.Background(), location)
	assert.NoError(t, err)
	assert.Equal(t, testImage, retrievedImage)
}
