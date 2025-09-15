package grpcConn

import (
	"context"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type ClientConnInterface interface {
	Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error
	NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error)
	Close() error
}

// MockGrpcClientConn is a mock implementation of grpc.ClientConn
type MockGrpcClientConn struct {
	mock.Mock
}

func (m *MockGrpcClientConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockGrpcClientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	return nil
}

func (m *MockGrpcClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, nil
}

func (m *MockGrpcClientConn) Target() string {
	return ""
}
