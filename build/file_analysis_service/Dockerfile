FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o file-analysis-service ./cmd/file_analysis_service

# Create a minimal image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/file-analysis-service .

# Create storage directory
RUN mkdir -p /app/storage/wordclouds

# Set environment variables
ENV PORT=50052
ENV DATABASE_URL="postgres://postgres:postgres@postgres:5432/textanalyzer?sslmode=disable"
ENV STORAGE_PATH="/app/storage/wordclouds"
ENV FILE_STORING_SERVICE_ADDRESS="file-storing-service:50051"
ENV WORDCLOUD_API_URL="https://quickchart.io/wordcloud"

# Expose the port
EXPOSE 50052

# Run the application
CMD ["./file-analysis-service"]