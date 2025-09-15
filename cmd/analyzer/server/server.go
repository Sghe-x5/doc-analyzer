package server

import (
	"context"
	"log"

	pb "local.dev/doc-analyzer/internal/proto/analyzer"
	"local.dev/doc-analyzer/internal/pkg/analyzer/service"
)

// Server implements the FileAnalysisServiceServer interface
type Server struct {
	pb.UnimplementedFileAnalysisServiceServer
	analysisService *service.AnalysisService
}

// NewServer creates a new Server instance
func NewServer(analysisService *service.AnalysisService) *Server {
	return &Server{
		analysisService: analysisService,
	}
}

// AnalyzeFile handles file analysis requests
func (s *Server) AnalyzeFile(ctx context.Context, req *pb.AnalyzeFileRequest) (*pb.AnalyzeFileResponse, error) {
	log.Printf("Received analysis request for file ID: %s", req.FileId)

	paragraphCount, wordCount, characterCount, isPlagiarism, similarFileIDs, wordCloudLocation, err := s.analysisService.AnalyzeFile(
		ctx,
		req.FileId,
		req.GenerateWordCloud,
	)
	if err != nil {
		log.Printf("Failed to analyze file: %v", err)
		return nil, err
	}

	log.Printf("File analyzed successfully: %s", req.FileId)
	return &pb.AnalyzeFileResponse{
		ParagraphCount:    paragraphCount,
		WordCount:         wordCount,
		CharacterCount:    characterCount,
		IsPlagiarism:      isPlagiarism,
		SimilarFileIds:    similarFileIDs,
		WordCloudLocation: wordCloudLocation,
	}, nil
}

// GetWordCloud handles word cloud retrieval requests
func (s *Server) GetWordCloud(ctx context.Context, req *pb.GetWordCloudRequest) (*pb.GetWordCloudResponse, error) {
	log.Printf("Received word cloud request for location: %s", req.Location)

	image, err := s.analysisService.GetWordCloud(ctx, req.Location)
	if err != nil {
		log.Printf("Failed to get word cloud: %v", err)
		return nil, err
	}

	log.Printf("Word cloud retrieved successfully: %s", req.Location)
	return &pb.GetWordCloudResponse{
		Image: image,
	}, nil
}