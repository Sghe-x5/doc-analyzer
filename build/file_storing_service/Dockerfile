FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o file-storing-service ./cmd/file_storing_service

# Create a minimal image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/file-storing-service .

# Create storage directory
RUN mkdir -p /app/storage/files

# Set environment variables
ENV PORT=50051
ENV DATABASE_URL="postgres://postgres:postgres@postgres:5432/textanalyzer?sslmode=disable"
ENV STORAGE_PATH="/app/storage/files"

# Expose the port
EXPOSE 50051

# Run the application
CMD ["./file-storing-service"]