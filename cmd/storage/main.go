package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"local.dev/doc-analyzer/cmd/storage/server"
	"local.dev/doc-analyzer/internal/pkg/storage/repository/postgres"
	"local.dev/doc-analyzer/internal/pkg/storage/service"
	"local.dev/doc-analyzer/internal/pkg/storage/storage/local"
	pb "local.dev/doc-analyzer/internal/proto/storage"
)

func main() {
	log.Println("Starting File Storing Service...")

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

	// Create the files table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			hash TEXT NOT NULL,
			location TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create files table: %v", err)
	}

	// Initialize storage
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage/files"
		log.Println("STORAGE_PATH not set, using default:", storagePath)
	}

	storage, err := local.NewLocalStorage(storagePath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize repository
	repo := postgres.NewFileRepo(db)

	// Initialize service
	fileService := service.NewFileService(repo, storage)

	// Initialize server
	grpcServer := grpc.NewServer()
	fileServer := server.NewServer(fileService)
	pb.RegisterFileStoringServiceServer(grpcServer, fileServer)

	// Start listening
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
		log.Println("PORT not set, using default:", port)
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("File Storing Service is listening on port %s...", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}