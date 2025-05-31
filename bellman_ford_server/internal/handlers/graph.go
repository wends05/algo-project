package handlers

import (
	"bellman_ford_server/internal/models"
	"bellman_ford_server/internal/services"
	"bellman_ford_server/internal/validators"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GraphHandler handles graph generation endpoints
type GraphHandler struct {
	service   *services.GraphService
	validator *validators.GraphValidator
}

// NewGraphHandler creates a new GraphHandler
func NewGraphHandler(service *services.GraphService, validator *validators.GraphValidator) *GraphHandler {
	return &GraphHandler{
		service:   service,
		validator: validator,
	}
}

// GenerateRandomGraph handles GET /graph
func (h *GraphHandler) GenerateRandomGraph(c *gin.Context) {
	// Parse query parameters
	params, err := h.parseGraphParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid parameters",
			Message: err.Error(),
		})
		return
	}

	// Validate parameters
	if err := h.validator.ValidateGenerationParams(params.Vertices, params.WeightRange, params.StartVertex, params.EndVertex); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
		return
	}

	// Log parameters
	fmt.Printf("Generating graph with %d vertices, weight range %d, start vertex %d, end vertex %d\n",
		params.Vertices, params.WeightRange, params.StartVertex, params.EndVertex)

	// Generate graph
	graph, err := h.service.GenerateRandomGraph(*params)
	if err != nil {
		if serviceErr, ok := err.(*services.InvalidGraphError); ok {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Invalid graph parameters",
				Message: serviceErr.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Internal server error",
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, graph)
}

// parseGraphParams parses and validates query parameters for graph generation
func (h *GraphHandler) parseGraphParams(c *gin.Context) (*models.GraphGenerationParams, error) {
	// Get number of vertices from query parameter
	verticesStr := c.Query("number_of_vertices")
	if verticesStr == "" {
		return nil, fmt.Errorf("number_of_vertices parameter is required")
	}

	vertices, err := strconv.Atoi(verticesStr)
	if err != nil {
		return nil, fmt.Errorf("vertices must be a valid integer")
	}

	// Get weight range from query parameter, default to 100
	weightRangeStr := c.Query("weight_range")
	if weightRangeStr == "" {
		weightRangeStr = "100"
	}

	weightRange, err := strconv.Atoi(weightRangeStr)
	if err != nil {
		return nil, fmt.Errorf("weight_range must be a valid integer")
	}

	// Parse start and end vertices
	var startVertex, endVertex int

	startVertexStr := c.Query("start_vertex")
	if startVertexStr != "" {
		startVertex, err = strconv.Atoi(startVertexStr)
		if err != nil {
			return nil, fmt.Errorf("start_vertex must be a valid integer")
		}
	} else {
		startVertex = 0 // Default to 0
	}

	endVertexStr := c.Query("end_vertex")
	if endVertexStr != "" {
		endVertex, err = strconv.Atoi(endVertexStr)
		if err != nil {
			return nil, fmt.Errorf("end_vertex must be a valid integer")
		}
	} else {
		endVertex = vertices - 1 // Default to last vertex
	}

	return &models.GraphGenerationParams{
		Vertices:    vertices,
		WeightRange: weightRange,
		StartVertex: startVertex,
		EndVertex:   endVertex,
	}, nil
}
