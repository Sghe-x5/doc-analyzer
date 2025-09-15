package clients_test

import (
	"context"
	"errors"
	"local.dev/doc-analyzer/internal/pkg/grpcConn"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"local.dev/doc-analyzer/internal/pkg/analyzer/clients"
	pb "local.dev/doc-analyzer/internal/proto/storage"
	"local.dev/doc-analyzer/internal/proto/storage/mocks"
)

func TestFileStoringClientImpl_GetFile(t *testing.T) {
	// Create mock client
	mockClient := new(mocks.MockFileStoringServiceClient)

	// Create client with mock
	client := clients.NewFileStoringClientWithClient(mockClient, nil)

	// Test case: successful get
	t.Run("Successful get", func(t *testing.T) {
		// Set up mock expectations
		mockClient.On("GetFile", mock.Anything, &pb.GetFileRequest{
			FileId: "file123",
		}).Return(&pb.GetFileResponse{
			FileName: "test.txt",
			Content:  []byte("test content"),
		}, nil).Once()

		// Call the method
		fileName, content, err := client.GetFile(context.Background(), "file123")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "test.txt", fileName)
		assert.Equal(t, []byte("test content"), content)

		mockClient.AssertExpectations(t)
	})

	// Test case: error from service
	t.Run("Error from service", func(t *testing.T) {
		// Set up mock expectations
		mockClient.On("GetFile", mock.Anything, &pb.GetFileRequest{
			FileId: "file456",
		}).Return(nil, errors.New("connection error")).Once()

		// Call the method
		_, _, err := client.GetFile(context.Background(), "file456")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get file")

		mockClient.AssertExpectations(t)
	})
}

func TestNewFileStoringClient(t *testing.T) {
	// Test case: invalid address
	t.Run("Invalid address", func(t *testing.T) {
		// Use an invalid address that will cause an error
		client, err := clients.NewFileStoringClient("invalid-address:12345")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "failed to connect")
	})

	// Test case: valid address but no server (connection refused)
	t.Run("Connection refused", func(t *testing.T) {
		// Use a valid address format but no server is running
		// Find an unused port
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("Failed to find unused port: %v", err)
		}
		addr := listener.Addr().String()
		listener.Close() // Close the listener to free the port

		// Try to connect to the port (should fail with connection refused)
		client, err := clients.NewFileStoringClient(addr)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

func TestFileStoringClientImpl_Close(t *testing.T) {
	// Test case: nil connection
	t.Run("Nil connection", func(t *testing.T) {
		// Create client with nil connection
		client := clients.NewFileStoringClientWithClient(nil, nil)

		// Call the method
		err := client.Close()

		// Assert
		assert.NoError(t, err)
	})

	// Test case: with connection
	t.Run("With connection", func(t *testing.T) {
		// Create mock connection
		mockConn := new(grpcConn.MockGrpcClientConn)
		mockConn.On("Close").Return(nil).Once()

		// Create client with mock connection
		client := clients.NewFileStoringClientWithClient(nil, mockConn)

		// Call the method
		err := client.Close()

		// Assert
		assert.NoError(t, err)
		mockConn.AssertExpectations(t)
	})

	// Test case: connection close error
	t.Run("Connection close error", func(t *testing.T) {
		// Create mock connection
		mockConn := new(grpcConn.MockGrpcClientConn)
		mockConn.On("Close").Return(errors.New("close error")).Once()

		// Create client with mock connection
		client := clients.NewFileStoringClientWithClient(nil, mockConn)

		// Call the method
		err := client.Close()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "close error", err.Error())
		mockConn.AssertExpectations(t)
	})
}
