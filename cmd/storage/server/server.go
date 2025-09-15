package server

import (
	"context"
	"log"

	pb "local.dev/doc-analyzer/internal/proto/storage"
	"local.dev/doc-analyzer/internal/pkg/storage/service"
)

// Server implements the FileStoringServiceServer interface
type Server struct {
	pb.UnimplementedFileStoringServiceServer
	fileService *service.FileService
}

// NewServer creates a new Server instance
func NewServer(fileService *service.FileService) *Server {
	return &Server{
		fileService: fileService,
	}
}

// UploadFile handles file upload requests
func (s *Server) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	log.Printf("Received upload request for file: %s", req.FileName)

	fileID, err := s.fileService.UploadFile(ctx, req.FileName, req.Content)
	if err != nil {
		log.Printf("Failed to upload file: %v", err)
		return nil, err
	}

	log.Printf("File uploaded successfully with ID: %s", fileID)
	return &pb.UploadFileResponse{
		FileId: fileID,
	}, nil
}

// GetFile handles file retrieval requests
func (s *Server) GetFile(ctx context.Context, req *pb.GetFileRequest) (*pb.GetFileResponse, error) {
	log.Printf("Received get file request for ID: %s", req.FileId)

	fileName, content, err := s.fileService.GetFile(ctx, req.FileId)
	if err != nil {
		log.Printf("Failed to get file: %v", err)
		return nil, err
	}

	log.Printf("File retrieved successfully: %s", fileName)
	return &pb.GetFileResponse{
		FileName: fileName,
		Content:  content,
	}, nil
}