package clients

import (
	"context"
	"fmt"
	"local.dev/doc-analyzer/internal/pkg/grpcConn"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "local.dev/doc-analyzer/internal/proto/storage"
)

// FileStoringClientInterface defines the interface for the File Storing Client
type FileStoringClientInterface interface {
	GetFile(ctx context.Context, fileID string) (string, []byte, error)
	Close() error
}

// FileStoringClient provides methods for interacting with the File Storing Service
type FileStoringClient struct {
	client pb.FileStoringServiceClient
	conn   grpcConn.ClientConnInterface
}

// NewFileStoringClient creates a new FileStoringClient instance
func NewFileStoringClient(address string) (*FileStoringClient, error) {
	// Set up a connection to the server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to File Storing Service: %w", err)
	}

	client := pb.NewFileStoringServiceClient(conn)

	return &FileStoringClient{
		client: client,
		conn:   conn,
	}, nil
}

// NewFileStoringClientWithClient creates a new FileStoringClient instance with a provided client
// This is primarily used for testing
func NewFileStoringClientWithClient(client pb.FileStoringServiceClient, conn grpcConn.ClientConnInterface) *FileStoringClient {
	return &FileStoringClient{
		client: client,
		conn:   conn,
	}
}

// Close closes the client connection
func (c *FileStoringClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetFile retrieves a file from the File Storing Service
func (c *FileStoringClient) GetFile(ctx context.Context, fileID string) (string, []byte, error) {
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Make the request
	resp, err := c.client.GetFile(ctx, &pb.GetFileRequest{
		FileId: fileID,
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to get file: %w", err)
	}

	return resp.FileName, resp.Content, nil
}
