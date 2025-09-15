package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	pb "local.dev/doc-analyzer/internal/proto/analyzer"
)

// FileAnalysisClientInterface defines the interface for the File Analysis Client
type FileAnalysisClientInterface interface {
	AnalyzeFile(ctx context.Context, fileID string, generateWordCloud bool) (*pb.AnalyzeFileResponse, error)
	GetWordCloud(ctx context.Context, location string) ([]byte, error)
	Close() error
}

// AnalysisHandler handles file analysis operations
type AnalysisHandler struct {
	client FileAnalysisClientInterface
}

// NewAnalysisHandler creates a new AnalysisHandler instance
func NewAnalysisHandler(client FileAnalysisClientInterface) *AnalysisHandler {
	return &AnalysisHandler{client: client}
}

// AnalyzeFileRequest represents the request body for file analysis
type AnalyzeFileRequest struct {
	FileID            string `json:"file_id" binding:"required" example:"file123"`
	GenerateWordCloud bool   `json:"generate_word_cloud" example:"true"`
}

// AnalyzeFileResponse represents the response for file analysis
type AnalyzeFileResponse struct {
	ParagraphCount    int32    `json:"paragraph_count" example:"5"`
	WordCount         int32    `json:"word_count" example:"100"`
	CharacterCount    int32    `json:"character_count" example:"500"`
	IsPlagiarism      bool     `json:"is_plagiarism" example:"false"`
	SimilarFileIds    []string `json:"similar_file_ids" example:"[]"`
	WordCloudLocation string   `json:"word_cloud_location" example:"wordclouds/file123.png"`
}

// AnalyzeFile godoc
// @Summary Analyze a file
// @Description Analyze a file by its ID
// @Tags analysis
// @Accept json
// @Produce json
// @Param request body AnalyzeFileRequest true "Analysis request"
// @Success 200 {object} AnalyzeFileResponse "Analysis results"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/analysis [post]
func (h *AnalysisHandler) AnalyzeFile(c *gin.Context) {
	var request AnalyzeFileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.AnalyzeFile(c.Request.Context(), request.FileID, request.GenerateWordCloud)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AnalyzeFileResponse{
		ParagraphCount:    resp.ParagraphCount,
		WordCount:         resp.WordCount,
		CharacterCount:    resp.CharacterCount,
		IsPlagiarism:      resp.IsPlagiarism,
		SimilarFileIds:    resp.SimilarFileIds,
		WordCloudLocation: resp.WordCloudLocation,
	})
}

// GetWordCloud godoc
// @Summary Get a word cloud
// @Description Get a word cloud image by its location
// @Tags analysis
// @Produce image/png
// @Param location path string true "Word cloud location"
// @Success 200 {file} binary "Word cloud image"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Word cloud not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/wordcloud/{location} [get]
func (h *AnalysisHandler) GetWordCloud(c *gin.Context) {
	location := c.Param("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location is required"})
		return
	}

	image, err := h.client.GetWordCloud(c.Request.Context(), location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/png", image)
}
