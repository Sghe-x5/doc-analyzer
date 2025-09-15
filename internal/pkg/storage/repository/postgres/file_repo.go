package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"local.dev/doc-analyzer/internal/pkg/storage/repository"
)

// FileRepo implements the FileRepository interface using PostgreSQL
type FileRepo struct {
	db *sql.DB
}

// NewFileRepo creates a new FileRepo instance
func NewFileRepo(db *sql.DB) repository.FileRepository {
	return &FileRepo{db: db}
}

// SaveFile saves file metadata to the database
func (r *FileRepo) SaveFile(ctx context.Context, id, name, hash, location string) error {
	query := `
		INSERT INTO files (id, name, hash, location, created_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
	`
	_, err := r.db.ExecContext(ctx, query, id, name, hash, location)
	if err != nil {
		return fmt.Errorf("failed to save file metadata: %w", err)
	}
	return nil
}

// GetFileByID retrieves file metadata by ID
func (r *FileRepo) GetFileByID(ctx context.Context, id string) (string, string, error) {
	query := `
		SELECT name, location FROM files WHERE id = $1
	`
	var name, location string
	err := r.db.QueryRowContext(ctx, query, id).Scan(&name, &location)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", fmt.Errorf("file not found with id %s", id)
		}
		return "", "", fmt.Errorf("failed to get file by id: %w", err)
	}
	return name, location, nil
}

// GetFileByHash retrieves file metadata by hash
func (r *FileRepo) GetFileByHash(ctx context.Context, hash string) (string, error) {
	query := `
		SELECT id FROM files WHERE hash = $1
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, hash).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil // No error, just no file with this hash
		}
		return "", fmt.Errorf("failed to get file by hash: %w", err)
	}
	return id, nil
}