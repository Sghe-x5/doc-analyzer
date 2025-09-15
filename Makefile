# Generate proto files
proto:
	bash scripts/generate_protos.sh

# Download Google proto files
proto-download:
	bash scripts/download_protos.sh

# Generate Swagger documentation
swagger:
	mkdir -p internal/pkg/api_gateway/docs
	swag init -g cmd/api_gateway/main.go -o internal/pkg/api_gateway/docs

# Run all tests and generate coverage report
test-coverage:
	mkdir -p coverage
	go test -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	go tool cover -func=coverage/coverage.out

# Calculate overall coverage excluding generated files
coverage-stats: test-coverage
	bash scripts/calculate_coverage.sh
