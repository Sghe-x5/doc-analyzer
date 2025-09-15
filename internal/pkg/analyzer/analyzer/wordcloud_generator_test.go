package analyzer_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"local.dev/doc-analyzer/internal/pkg/analyzer/analyzer"
)

func TestWordCloudGenerator_NewWordCloudGenerator(t *testing.T) {
	// Test with custom URL
	customURL := "https://example.com/wordcloud"
	generator := analyzer.NewWordCloudGenerator(customURL)
	assert.NotNil(t, generator, "Generator should not be nil")

	// Test with empty URL (should use default)
	generator = analyzer.NewWordCloudGenerator("")
	assert.NotNil(t, generator, "Generator should not be nil")
}

func TestWordCloudGenerator_GenerateWordCloud(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		assert.Equal(t, "POST", r.Method, "Request method should be POST")
		
		// Check content type
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"), "Content-Type should be application/json")
		
		// Return a mock image
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mock-image-data"))
	}))
	defer server.Close()

	// Create a generator with the mock server URL
	generator := analyzer.NewWordCloudGenerator(server.URL)

	// Test generating a word cloud
	imageData, location, err := generator.GenerateWordCloud(context.Background(), "This is a test text for word cloud generation")
	
	// Assertions
	assert.NoError(t, err, "Should not return an error")
	assert.Equal(t, []byte("mock-image-data"), imageData, "Image data should match")
	assert.NotEmpty(t, location, "Location should not be empty")
	assert.Contains(t, location, ".png", "Location should be a PNG file")
}

func TestWordCloudGenerator_GenerateWordCloud_ServerError(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	// Create a generator with the mock server URL
	generator := analyzer.NewWordCloudGenerator(server.URL)

	// Test generating a word cloud with server error
	imageData, location, err := generator.GenerateWordCloud(context.Background(), "This is a test text")
	
	// Assertions
	assert.Error(t, err, "Should return an error")
	assert.Nil(t, imageData, "Image data should be nil")
	assert.Empty(t, location, "Location should be empty")
}

func TestWordCloudGenerator_GenerateWordCloud_InvalidURL(t *testing.T) {
	// Create a generator with an invalid URL
	generator := analyzer.NewWordCloudGenerator("http://invalid-url-that-does-not-exist.example")

	// Test generating a word cloud with an invalid URL
	imageData, location, err := generator.GenerateWordCloud(context.Background(), "This is a test text")
	
	// Assertions
	assert.Error(t, err, "Should return an error")
	assert.Nil(t, imageData, "Image data should be nil")
	assert.Empty(t, location, "Location should be empty")
}

func TestWordCloudGenerator_GenerateWordCloud_CanceledContext(t *testing.T) {
	// Create a generator
	generator := analyzer.NewWordCloudGenerator("")

	// Create a canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	// Test generating a word cloud with a canceled context
	imageData, location, err := generator.GenerateWordCloud(ctx, "This is a test text")
	
	// Assertions
	assert.Error(t, err, "Should return an error")
	assert.Nil(t, imageData, "Image data should be nil")
	assert.Empty(t, location, "Location should be empty")
}