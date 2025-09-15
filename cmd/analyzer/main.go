package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"local.dev/doc-analyzer/cmd/analyzer/server"
	"local.dev/doc-analyzer/internal/pkg/analyzer/analyzer"
	"local.dev/doc-analyzer/internal/pkg/analyzer/clients"
	"local.dev/doc-analyzer/internal/pkg/analyzer/repository/postgres"
	"local.dev/doc-analyzer/internal/pkg/analyzer/service"
	"local.dev/doc-analyzer/internal/pkg/analyzer/storage/local"
	pb "local.dev/doc-analyzer/internal/proto/analyzer"
)

func main() {
	log.Println("Starting File Analysis Service...")

	// Get database connection string from environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@postgres:5432/textanalyzer?sslmode=disable"
		log.Println("DATABASE_URL not set, using default:", dbURL)
	}

	// Connect to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to the database")

	// Create the analysis_results and similar_files tables if they don't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS analysis_results (
			file_id TEXT PRIMARY KEY,
			paragraph_count INT NOT NULL,
			word_count INT NOT NULL,
			character_count INT NOT NULL,
			is_plagiarism BOOLEAN NOT NULL,
			word_cloud_location TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE TABLE IF NOT EXISTS similar_files (
			file_id TEXT,
			similar_file_id TEXT,
			PRIMARY KEY (file_id, similar_file_id)
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Initialize storage
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage/wordclouds"
		log.Println("STORAGE_PATH not set, using default:", storagePath)
	}

	storage, err := local.NewLocalStorage(storagePath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize repository
	repo := postgres.NewAnalysisRepo(db)

	// Initialize File Storing Service client
	fileStoringAddress := os.Getenv("FILE_STORING_SERVICE_ADDRESS")
	if fileStoringAddress == "" {
		fileStoringAddress = "file-storing-service:50051"
		log.Println("FILE_STORING_SERVICE_ADDRESS not set, using default:", fileStoringAddress)
	}

	fileStoringClient, err := clients.NewFileStoringClient(fileStoringAddress)
	if err != nil {
		log.Fatalf("Failed to initialize File Storing Service client: %v", err)
	}
	defer fileStoringClient.Close()

	// Initialize analyzers
	textAnalyzer := analyzer.NewTextAnalyzer()
	plagiarismChecker := analyzer.NewPlagiarismChecker()
	
	wordCloudAPIURL := os.Getenv("WORDCLOUD_API_URL")
	if wordCloudAPIURL == "" {
		wordCloudAPIURL = "https://quickchart.io/wordcloud"
		log.Println("WORDCLOUD_API_URL not set, using default:", wordCloudAPIURL)
	}
	wordCloudGenerator := analyzer.NewWordCloudGenerator(wordCloudAPIURL)

	// Initialize service
	analysisService := service.NewAnalysisService(
		repo,
		storage,
		fileStoringClient,
		textAnalyzer,
		plagiarismChecker,
		wordCloudGenerator,
	)

	// Initialize server
	grpcServer := grpc.NewServer()
	analysisServer := server.NewServer(analysisService)
	pb.RegisterFileAnalysisServiceServer(grpcServer, analysisServer)

	// Start listening
	port := os.Getenv("PORT")
	if port == "" {
		port = "50052"
		log.Println("PORT not set, using default:", port)
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("File Analysis Service is listening on port %s...", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}