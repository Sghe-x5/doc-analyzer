package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"

	"local.dev/doc-analyzer/internal/pkg/storage/repository"
	"local.dev/doc-analyzer/internal/pkg/storage/storage"
)

// FileService handles the business logic for file operations
type FileService struct {
	repo    repository.FileRepository
	storage storage.FileStorage
}

// NewFileService creates a new FileService instance
func NewFileService(repo repository.FileRepository, storage storage.FileStorage) *FileService {
	return &FileService{
		repo:    repo,
		storage: storage,
	}
}

// UploadFile handles the file upload process
func (s *FileService) UploadFile(ctx context.Context, fileName string, content []byte) (string, error) {
	// Calculate file hash
	hash := sha256.Sum256(content)
	hashStr := hex.EncodeToString(hash[:])

	// Check if file with this hash already exists
	fileID, err := s.repo.GetFileByHash(ctx, hashStr)
	if err != nil {
		return "", fmt.Errorf("failed to check file existence: %w", err)
	}

	// If file exists, return its ID
	if fileID != "" {
		return fileID, nil
	}

	// Generate a new file ID
	fileID = uuid.New().String()

	// Define file location
	location := fileID // Using fileID as location for simplicity

	// Save file content to storage
	if err := s.storage.SaveFile(ctx, location, content); err != nil {
		return "", fmt.Errorf("failed to save file content: %w", err)
	}

	// Save file metadata to repository
	if err := s.repo.SaveFile(ctx, fileID, fileName, hashStr, location); err != nil {
		return "", fmt.Errorf("failed to save file metadata: %w", err)
	}

	return fileID, nil
}

// GetFile retrieves a file by its ID
func (s *FileService) GetFile(ctx context.Context, fileID string) (string, []byte, error) {
	// Get file metadata from repository
	fileName, location, err := s.repo.GetFileByID(ctx, fileID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get file metadata: %w", err)
	}

	// Get file content from storage
	content, err := s.storage.GetFile(ctx, location)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get file content: %w", err)
	}

	return fileName, content, nil
}