package service

import (
	"context"
	"fmt"

	"local.dev/doc-analyzer/internal/pkg/analyzer/analyzer"
	"local.dev/doc-analyzer/internal/pkg/analyzer/clients"
	"local.dev/doc-analyzer/internal/pkg/analyzer/repository"
	"local.dev/doc-analyzer/internal/pkg/analyzer/storage"
)

// AnalysisService handles the business logic for file analysis operations
type AnalysisService struct {
	repo               repository.AnalysisRepository
	storage            storage.WordCloudStorage
	fileStoringClient  clients.FileStoringClientInterface
	textAnalyzer       *analyzer.TextAnalyzer
	plagiarismChecker  *analyzer.PlagiarismChecker
	wordCloudGenerator *analyzer.WordCloudGenerator
}

// NewAnalysisService creates a new AnalysisService instance
func NewAnalysisService(
	repo repository.AnalysisRepository,
	storage storage.WordCloudStorage,
	fileStoringClient clients.FileStoringClientInterface,
	textAnalyzer *analyzer.TextAnalyzer,
	plagiarismChecker *analyzer.PlagiarismChecker,
	wordCloudGenerator *analyzer.WordCloudGenerator,
) *AnalysisService {
	return &AnalysisService{
		repo:               repo,
		storage:            storage,
		fileStoringClient:  fileStoringClient,
		textAnalyzer:       textAnalyzer,
		plagiarismChecker:  plagiarismChecker,
		wordCloudGenerator: wordCloudGenerator,
	}
}

// AnalyzeFile analyzes a file and returns the analysis results
func (s *AnalysisService) AnalyzeFile(ctx context.Context, fileID string, generateWordCloud bool) (
	paragraphCount, wordCount, characterCount int32,
	isPlagiarism bool,
	similarFileIDs []string,
	wordCloudLocation string,
	err error,
) {
	// Try to get existing analysis results
	paragraphCount, wordCount, characterCount, isPlagiarism, wordCloudLocation, err = s.repo.GetAnalysisResult(ctx, fileID)
	if err == nil {
		// Analysis results exist, get similar file IDs if it's plagiarism
		if isPlagiarism {
			similarFileIDs, err = s.repo.GetSimilarFiles(ctx, fileID)
			if err != nil {
				return 0, 0, 0, false, nil, "", fmt.Errorf("failed to get similar files: %w", err)
			}
		}
		return paragraphCount, wordCount, characterCount, isPlagiarism, similarFileIDs, wordCloudLocation, nil
	}

	// Get file content from File Storing Service
	_, content, err := s.fileStoringClient.GetFile(ctx, fileID)
	if err != nil {
		return 0, 0, 0, false, nil, "", fmt.Errorf("failed to get file content: %w", err)
	}

	// Convert content to string
	contentStr := string(content)

	// Analyze text
	paragraphCount, wordCount, characterCount = s.textAnalyzer.AnalyzeText(contentStr)

	// Check for plagiarism
	// First, get all other file IDs
	otherFileIDs, err := s.repo.GetAllFileIDs(ctx)
	if err != nil {
		return 0, 0, 0, false, nil, "", fmt.Errorf("failed to get all file IDs: %w", err)
	}

	// Get content of all other files
	otherContents := make(map[string]string)
	for _, otherFileID := range otherFileIDs {
		if otherFileID == fileID {
			continue // Skip the current file
		}

		_, otherContent, err := s.fileStoringClient.GetFile(ctx, otherFileID)
		if err != nil {
			// Log the error but continue with other files
			fmt.Printf("Failed to get content for file %s: %v\n", otherFileID, err)
			continue
		}

		otherContents[otherFileID] = string(otherContent)
	}

	// Check for plagiarism
	isPlagiarism, similarFileIDs = s.plagiarismChecker.CheckPlagiarism(ctx, contentStr, otherContents)

	// Generate word cloud if requested
	if generateWordCloud {
		var text []byte
		_, text, err = s.fileStoringClient.GetFile(ctx, fileID)

		if err != nil {
			return 0, 0, 0, false, nil, "", fmt.Errorf("failed to get file content: %w", err)
		}

		// Generate word cloud
		wordCloudImage, location, err := s.wordCloudGenerator.GenerateWordCloud(ctx, string(text))
		if err != nil {
			// Log the error but continue without word cloud
			fmt.Printf("Failed to generate word cloud: %v\n", err)
		} else {
			// Save word cloud image
			err = s.storage.SaveWordCloud(ctx, location, wordCloudImage)
			if err != nil {
				fmt.Printf("Failed to save word cloud: %v\n", err)
			} else {
				wordCloudLocation = location
			}
		}
	}

	// Save analysis results
	err = s.repo.SaveAnalysisResult(ctx, fileID, paragraphCount, wordCount, characterCount, isPlagiarism, wordCloudLocation)
	if err != nil {
		return 0, 0, 0, false, nil, "", fmt.Errorf("failed to save analysis results: %w", err)
	}

	// Save similar files if plagiarism is detected
	if isPlagiarism {
		for _, similarFileID := range similarFileIDs {
			err = s.repo.SaveSimilarFile(ctx, fileID, similarFileID)
			if err != nil {
				// Log the error but continue with other similar files
				fmt.Printf("Failed to save similar file %s: %v\n", similarFileID, err)
			}
		}
	}

	return paragraphCount, wordCount, characterCount, isPlagiarism, similarFileIDs, wordCloudLocation, nil
}

// GetWordCloud retrieves a word cloud image by its location
func (s *AnalysisService) GetWordCloud(ctx context.Context, location string) ([]byte, error) {
	return s.storage.GetWordCloud(ctx, location)
}
