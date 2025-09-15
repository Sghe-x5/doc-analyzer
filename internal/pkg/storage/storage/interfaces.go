package storage

import (
	"context"
)

// FileStorage defines the interface for file content operations
type FileStorage interface {
	// SaveFile saves file content to storage
	SaveFile(ctx context.Context, location string, content []byte) error
	
	// GetFile retrieves file content from storage
	GetFile(ctx context.Context, location string) ([]byte, error)
}