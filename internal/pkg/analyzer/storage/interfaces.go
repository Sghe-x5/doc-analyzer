package storage

import (
	"context"
)

// WordCloudStorage defines the interface for word cloud image operations
type WordCloudStorage interface {
	// SaveWordCloud saves a word cloud image to storage
	SaveWordCloud(ctx context.Context, location string, image []byte) error
	
	// GetWordCloud retrieves a word cloud image from storage
	GetWordCloud(ctx context.Context, location string) ([]byte, error)
}