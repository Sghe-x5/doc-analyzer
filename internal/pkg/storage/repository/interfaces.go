package repository

import (
	"context"
)

// FileRepository defines the interface for file metadata operations
type FileRepository interface {
	// SaveFile saves file metadata to the database
	SaveFile(ctx context.Context, id, name, hash, location string) error
	
	// GetFileByID retrieves file metadata by ID
	GetFileByID(ctx context.Context, id string) (name string, location string, err error)
	
	// GetFileByHash retrieves file metadata by hash
	GetFileByHash(ctx context.Context, hash string) (id string, err error)
}