#!/bin/bash

# Create directories for generated code
mkdir -p ./internal/proto/api_gateway
mkdir -p ./internal/proto/file_analysis_service
mkdir -p ./internal/proto/file_storing_service

# Create directory for Swagger documentation
mkdir -p ./api/swagger

# Generate code for API Gateway
protoc -I ./proto \
   --go_out ./internal/proto/api_gateway --go_opt paths=source_relative \
   --go-grpc_out ./internal/proto/api_gateway --go-grpc_opt paths=source_relative \
   --grpc-gateway_out ./internal/proto/api_gateway --grpc-gateway_opt paths=source_relative \
   --openapiv2_out ./api/swagger --openapiv2_opt logtostderr=true \
   proto/api_gateway.proto

# Generate code for File Analysis Service
protoc -I ./proto \
   --go_out ./internal/proto/file_analysis_service --go_opt paths=source_relative \
   --go-grpc_out ./internal/proto/file_analysis_service --go-grpc_opt paths=source_relative \
   --grpc-gateway_out ./internal/proto/file_analysis_service --grpc-gateway_opt paths=source_relative \
   --openapiv2_out ./api/swagger --openapiv2_opt logtostderr=true \
   proto/file_analysis_service.proto

# Generate code for File Storing Service
protoc -I ./proto \
   --go_out ./internal/proto/file_storing_service --go_opt paths=source_relative \
   --go-grpc_out ./internal/proto/file_storing_service --go-grpc_opt paths=source_relative \
   --grpc-gateway_out ./internal/proto/file_storing_service --grpc-gateway_opt paths=source_relative \
   --openapiv2_out ./api/swagger --openapiv2_opt logtostderr=true \
   proto/file_storing_service.proto

# Merge Swagger files into a single file
echo "Merging Swagger files..."
jq -s 'reduce .[] as $item ({}; . * $item)' ./api/swagger/*.swagger.json > ./api/swagger/api.swagger.json
