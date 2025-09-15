FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./cmd/api_gateway

# Create a minimal image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/api-gateway .

# Set environment variables
ENV HTTP_PORT=8080
ENV FILE_STORING_SERVICE_ADDRESS="file-storing-service:50051"
ENV FILE_ANALYSIS_SERVICE_ADDRESS="file-analysis-service:50052"
ENV GIN_MODE=release

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./api-gateway"]
