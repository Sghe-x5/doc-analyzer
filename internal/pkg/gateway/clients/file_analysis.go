package clients

import (
	"context"
	"fmt"
	"local.dev/doc-analyzer/internal/pkg/grpcConn"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "local.dev/doc-analyzer/internal/proto/analyzer"
)

// FileAnalysisClient provides methods for interacting with the File Analysis Service
type FileAnalysisClient struct {
	client pb.FileAnalysisServiceClient
	conn   grpcConn.ClientConnInterface
}

// NewFileAnalysisClient creates a new FileAnalysisClient instance
func NewFileAnalysisClient(address string) (*FileAnalysisClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to File Analysis Service: %w", err)
	}

	client := pb.NewFileAnalysisServiceClient(conn)

	return &FileAnalysisClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the client connection
func (c *FileAnalysisClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// AnalyzeFile sends a request to analyze a file
func (c *FileAnalysisClient) AnalyzeFile(ctx context.Context, fileID string, generateWordCloud bool) (*pb.AnalyzeFileResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second) // Analysis might take longer
	defer cancel()

	maxRetries := 3
	retryDelay := 1 * time.Second

	var resp *pb.AnalyzeFileResponse
	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = c.client.AnalyzeFile(ctx, &pb.AnalyzeFileRequest{
			FileId:            fileID,
			GenerateWordCloud: generateWordCloud,
		})

		if err == nil {
			break
		}

		s, ok := status.FromError(err)
		if !ok || (s.Code() != codes.Unavailable && s.Code() != codes.DeadlineExceeded) {
			return nil, fmt.Errorf("failed to analyze file: %w", err)
		}

		if attempt == maxRetries-1 {
			return nil, fmt.Errorf("failed to analyze file after %d attempts: %w", maxRetries, err)
		}

		time.Sleep(retryDelay)
		retryDelay *= 2
	}

	return resp, nil
}

// GetWordCloud retrieves a word cloud image
func (c *FileAnalysisClient) GetWordCloud(ctx context.Context, location string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	maxRetries := 3
	retryDelay := 1 * time.Second

	var resp *pb.GetWordCloudResponse
	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = c.client.GetWordCloud(ctx, &pb.GetWordCloudRequest{
			Location: location,
		})

		if err == nil {
			break
		}

		s, ok := status.FromError(err)
		if !ok || (s.Code() != codes.Unavailable && s.Code() != codes.DeadlineExceeded) {
			return nil, fmt.Errorf("failed to get word cloud: %w", err)
		}

		if attempt == maxRetries-1 {
			return nil, fmt.Errorf("failed to get word cloud after %d attempts: %w", maxRetries, err)
		}

		time.Sleep(retryDelay)
		retryDelay *= 2
	}

	return resp.Image, nil
}
