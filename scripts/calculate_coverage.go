package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Open the coverage file
	file, err := os.Open("coverage/coverage.out")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening coverage file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Variables to track coverage
	var totalLines, coveredLines int

	// Read the file line by line
	scanner := bufio.NewScanner(file)

	// Skip the first line (mode: set)
	if scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "mode:") {
			fmt.Fprintf(os.Stderr, "Invalid coverage file format\n")
			os.Exit(1)
		}
	}

	// Process each line
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}

		filePath := parts[0]

		// Skip generated files and this script itself
		if isGeneratedFile(filePath) || isSelf(filePath) {
			continue
		}

		// Parse the line range and coverage information
		// Format is typically: file:startLine.startCol,endLine.endCol statements count
		lineRange := strings.Split(parts[1], " ")[0]
		lineNumbers := strings.Split(lineRange, ",")
		if len(lineNumbers) < 2 {
			continue
		}

		// Extract start and end line numbers
		startLineStr := strings.Split(lineNumbers[0], ".")[0]
		endLineStr := strings.Split(lineNumbers[1], ".")[0]

		startLine, err := strconv.Atoi(startLineStr)
		if err != nil {
			continue
		}

		endLine, err := strconv.Atoi(endLineStr)
		if err != nil {
			continue
		}

		// Calculate number of lines in this block
		linesInBlock := endLine - startLine + 1

		// Extract coverage information
		coverageInfo := strings.Split(parts[len(parts)-1], " ")
		if len(coverageInfo) < 3 {
			continue
		}

		// The last number indicates if the line is covered (0 = not covered, > 0 = covered)
		covered := coverageInfo[len(coverageInfo)-1] != "0"

		if covered {
			coveredLines += linesInBlock
		}
		totalLines += linesInBlock
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading coverage file: %v\n", err)
		os.Exit(1)
	}

	// Calculate coverage percentage
	var coveragePercent float64
	if totalLines > 0 {
		coveragePercent = float64(coveredLines) / float64(totalLines) * 100
	}

	fmt.Printf("Total lines (excluding generated files): %d\n", totalLines)
	fmt.Printf("Covered lines: %d\n", coveredLines)
	fmt.Printf("Coverage: %.1f%%\n", coveragePercent)
}

// isGeneratedFile returns true if the file is likely generated
func isGeneratedFile(filePath string) bool {
	// Check if the file is in the proto directory
	if strings.Contains(filePath, "/proto/") {
		return true
	}

	// Check if the file is a generated Swagger doc
	if strings.Contains(filePath, "/gateway/docs/") {
		return true
	}

	// Check if the file is a generated Swagger doc
	if strings.Contains(filePath, "/mocks/") {
		return true
	}

	// Check file extensions commonly associated with generated files
	generatedExtensions := []string{".pb.go", ".pb.gw.go"}
	for _, genExt := range generatedExtensions {
		if strings.HasSuffix(filePath, genExt) {
			return true
		}
	}

	return false
}

// isSelf returns true if the file is this script itself
func isSelf(filePath string) bool {
	// Check if the file is this script
	return strings.HasSuffix(filePath, "scripts/calculate_coverage.go")
}
