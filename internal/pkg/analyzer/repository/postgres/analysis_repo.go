package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"local.dev/doc-analyzer/internal/pkg/analyzer/repository"
)

// AnalysisRepo implements the AnalysisRepository interface using PostgreSQL
type AnalysisRepo struct {
	db *sql.DB
}

// NewAnalysisRepo creates a new AnalysisRepo instance
func NewAnalysisRepo(db *sql.DB) repository.AnalysisRepository {
	return &AnalysisRepo{db: db}
}

// SaveAnalysisResult saves analysis results to the database
func (r *AnalysisRepo) SaveAnalysisResult(ctx context.Context, fileID string, paragraphCount, wordCount, characterCount int32, isPlagiarism bool, wordCloudLocation string) error {
	query := `
		INSERT INTO analysis_results (
			file_id, paragraph_count, word_count, character_count, 
			is_plagiarism, word_cloud_location, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
		ON CONFLICT (file_id) DO UPDATE SET
			paragraph_count = $2,
			word_count = $3,
			character_count = $4,
			is_plagiarism = $5,
			word_cloud_location = $6,
			created_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.ExecContext(
		ctx, query, fileID, paragraphCount, wordCount, characterCount,
		isPlagiarism, wordCloudLocation,
	)
	if err != nil {
		return fmt.Errorf("failed to save analysis result: %w", err)
	}
	return nil
}

// GetAnalysisResult retrieves analysis results by file ID
func (r *AnalysisRepo) GetAnalysisResult(ctx context.Context, fileID string) (int32, int32, int32, bool, string, error) {
	query := `
		SELECT paragraph_count, word_count, character_count, is_plagiarism, word_cloud_location
		FROM analysis_results
		WHERE file_id = $1
	`
	var paragraphCount, wordCount, characterCount int32
	var isPlagiarism bool
	var wordCloudLocation sql.NullString

	err := r.db.QueryRowContext(ctx, query, fileID).Scan(
		&paragraphCount, &wordCount, &characterCount, &isPlagiarism, &wordCloudLocation,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, 0, false, "", fmt.Errorf("analysis result not found for file ID %s", fileID)
		}
		return 0, 0, 0, false, "", fmt.Errorf("failed to get analysis result: %w", err)
	}

	var location string
	if wordCloudLocation.Valid {
		location = wordCloudLocation.String
	}

	return paragraphCount, wordCount, characterCount, isPlagiarism, location, nil
}

// SaveSimilarFile saves information about a similar file (for plagiarism detection)
func (r *AnalysisRepo) SaveSimilarFile(ctx context.Context, fileID, similarFileID string) error {
	query := `
		INSERT INTO similar_files (file_id, similar_file_id)
		VALUES ($1, $2)
		ON CONFLICT (file_id, similar_file_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, fileID, similarFileID)
	if err != nil {
		return fmt.Errorf("failed to save similar file: %w", err)
	}
	return nil
}

// GetSimilarFiles retrieves IDs of similar files for a given file ID
func (r *AnalysisRepo) GetSimilarFiles(ctx context.Context, fileID string) ([]string, error) {
	query := `
		SELECT similar_file_id FROM similar_files WHERE file_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to query similar files: %w", err)
	}
	defer rows.Close()

	var similarFileIDs []string
	for rows.Next() {
		var similarFileID string
		if err := rows.Scan(&similarFileID); err != nil {
			return nil, fmt.Errorf("failed to scan similar file ID: %w", err)
		}
		similarFileIDs = append(similarFileIDs, similarFileID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over similar files: %w", err)
	}

	return similarFileIDs, nil
}

// GetAllFileIDs retrieves all file IDs in the database
func (r *AnalysisRepo) GetAllFileIDs(ctx context.Context) ([]string, error) {
	query := `
		SELECT file_id FROM analysis_results
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all file IDs: %w", err)
	}
	defer rows.Close()

	var fileIDs []string
	for rows.Next() {
		var fileID string
		if err := rows.Scan(&fileID); err != nil {
			return nil, fmt.Errorf("failed to scan file ID: %w", err)
		}
		fileIDs = append(fileIDs, fileID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over file IDs: %w", err)
	}

	return fileIDs, nil
}
