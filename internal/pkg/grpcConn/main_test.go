package grpcConn_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"local.dev/doc-analyzer/internal/pkg/grpcConn"
)

func TestMockGrpcClientConn_Close(t *testing.T) {
	// Create a new mock
	mockConn := &grpcConn.MockGrpcClientConn{}

	// Test case: Close returns nil
	t.Run("Close returns nil", func(t *testing.T) {
		// Setup expectations
		mockConn.On("Close").Return(nil).Once()

		// Call the method
		err := mockConn.Close()

		// Assert
		assert.NoError(t, err)
		mockConn.AssertExpectations(t)
	})

	// Test case: Close returns error
	t.Run("Close returns error", func(t *testing.T) {
		// Setup expectations
		expectedErr := errors.New("close error")
		mockConn.On("Close").Return(expectedErr).Once()

		// Call the method
		err := mockConn.Close()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockConn.AssertExpectations(t)
	})
}

func TestMockGrpcClientConn_Invoke(t *testing.T) {
	// Create a new mock
	mockConn := &grpcConn.MockGrpcClientConn{}

	// Test Invoke method
	t.Run("Invoke returns nil", func(t *testing.T) {
		// Call the method
		err := mockConn.Invoke(context.Background(), "test-method", "args", "reply")

		// Assert
		assert.NoError(t, err)
	})
}

func TestMockGrpcClientConn_NewStream(t *testing.T) {
	// Create a new mock
	mockConn := &grpcConn.MockGrpcClientConn{}

	// Test NewStream method
	t.Run("NewStream returns nil, nil", func(t *testing.T) {
		// Call the method
		stream, err := mockConn.NewStream(context.Background(), &grpc.StreamDesc{}, "test-method")

		// Assert
		assert.Nil(t, stream)
		assert.Nil(t, err)
	})
}

func TestMockGrpcClientConn_Target(t *testing.T) {
	// Create a new mock
	mockConn := &grpcConn.MockGrpcClientConn{}

	// Test Target method
	t.Run("Target returns empty string", func(t *testing.T) {
		// Call the method
		target := mockConn.Target()

		// Assert
		assert.Equal(t, "", target)
	})
}
