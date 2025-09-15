package repository

import (
	"context"
)

// AnalysisRepository defines the interface for analysis results operations
type AnalysisRepository interface {
	// SaveAnalysisResult saves analysis results to the database
	SaveAnalysisResult(ctx context.Context, fileID string, paragraphCount, wordCount, characterCount int32, isPlagiarism bool, wordCloudLocation string) error
	
	// GetAnalysisResult retrieves analysis results by file ID
	GetAnalysisResult(ctx context.Context, fileID string) (paragraphCount, wordCount, characterCount int32, isPlagiarism bool, wordCloudLocation string, err error)
	
	// SaveSimilarFile saves information about a similar file (for plagiarism detection)
	SaveSimilarFile(ctx context.Context, fileID, similarFileID string) error
	
	// GetSimilarFiles retrieves IDs of similar files for a given file ID
	GetSimilarFiles(ctx context.Context, fileID string) ([]string, error)
	
	// GetAllFileIDs retrieves all file IDs in the database
	GetAllFileIDs(ctx context.Context) ([]string, error)
}