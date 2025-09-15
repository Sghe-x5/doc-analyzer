package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// MockWordCloudGenerator is a mock implementation of the WordCloudGenerator
type MockWordCloudGenerator struct {
	mock.Mock
	apiURL string
}

// GenerateWordCloud mocks the GenerateWordCloud method
func (m *MockWordCloudGenerator) GenerateWordCloud(ctx context.Context, text string) ([]byte, string, error) {
	args := m.Called(ctx, text)
	if args.Get(0) == nil {
		return nil, args.String(1), args.Error(2)
	}
	return args.Get(0).([]byte), args.String(1), args.Error(2)
}