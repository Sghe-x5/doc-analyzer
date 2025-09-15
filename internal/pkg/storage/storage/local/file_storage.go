package local

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"local.dev/doc-analyzer/internal/pkg/storage/storage"
)

// LocalStorage implements the FileStorage interface using the local filesystem
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new LocalStorage instance
func NewLocalStorage(basePath string) (storage.FileStorage, error) {
	if basePath == "" {
		basePath = "."
	}

	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &LocalStorage{basePath: basePath}, nil
}

// SaveFile saves file content to the local filesystem
func (s *LocalStorage) SaveFile(ctx context.Context, location string, content []byte) error {
	if location == "" {
		location = "default.txt"
	}

	fullPath := filepath.Join(s.basePath, location)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// GetFile retrieves file content from the local filesystem
func (s *LocalStorage) GetFile(ctx context.Context, location string) ([]byte, error) {
	if location == "" {
		location = "default.txt"
	}

	fullPath := filepath.Join(s.basePath, location)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found at location %s", location)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return content, nil
}
