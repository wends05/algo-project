package handlers

import (
	"bellman_ford_server/internal/models"
	"bellman_ford_server/internal/services"
	"bellman_ford_server/internal/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BellmanFordHandler handles Bellman-Ford algorithm endpoints
type BellmanFordHandler struct {
	service   *services.BellmanFordService
	validator *validators.BellmanFordValidator
}

// NewBellmanFordHandler creates a new BellmanFordHandler
func NewBellmanFordHandler(service *services.BellmanFordService, validator *validators.BellmanFordValidator) *BellmanFordHandler {
	return &BellmanFordHandler{
		service:   service,
		validator: validator,
	}
}

// FindShortestPath handles POST /bellmanford
func (h *BellmanFordHandler) FindShortestPath(c *gin.Context) {
	var input models.BellmanFordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid input",
			Message: err.Error(),
		})
		return
	}

	// Validate input
	if err := h.validator.ValidateInput(input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
		return
	}

	// Process the request
	result, err := h.service.FindShortestPath(input)
	if err != nil {
		switch err.(type) {
		case *services.PathNotFoundError:
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "No path found",
				Message: err.Error(),
			})
		case *services.PathReconstructionError:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Path reconstruction error",
				Message: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Internal server error",
				Message: err.Error(),
			})
		}
		return
	}

	// Handle negative cycle case
	if result.HasNegativeCycle {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Negative cycle detected",
			Message: "The graph contains a negative cycle, making pathfinding impossible.",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
