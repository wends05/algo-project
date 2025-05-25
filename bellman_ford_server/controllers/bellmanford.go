package controllers

import (
	"bellman_ford_server/structs"
	"bellman_ford_server/utils"
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BellmanFord input
// BellmanFordInput represents the input data needed to run the Bellman-Ford algorithm
type BellmanFordInput struct {
	NumberOfVertices int            `json:"number_of_vertices"`
	Edges            []structs.Edge `json:"edges" binding:"dive"`
	StartVertex      int            `json:"start_vertex" binding:"min=0"`
	EndVertex        int            `json:"end_vertex" binding:"min=0"`
}

type BellmanFordResponse struct {
	Distance         int   `json:"distance"`
	Path             []int `json:"path"`
	HasNegativeCycle bool  `json:"has_negative_cycle" `
}

func BellmanFord(c *gin.Context) {

	data := BellmanFordInput{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Error:   "Invalid input",
			Message: err.Error(),
		})
		return
	}

	// --- Input Validation (more comprehensive) ---

	// Check if the number of vertices is valid
	if data.NumberOfVertices < 2 {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Error:   "Invalid graph size",
			Message: "Number of vertices must be at least 2 for pathfinding between distinct points.",
		})
		return
	}
	// Check if the number of edges is valid
	if data.StartVertex < 0 || data.StartVertex >= data.NumberOfVertices {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{Error: "Invalid start vertex"})
		return
	}
	// Check if the start vertex is valid
	if data.EndVertex < 0 || data.EndVertex >= data.NumberOfVertices {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{Error: "Invalid end vertex"})
		return
	}
	// Check if the start and end vertices are the same
	if data.StartVertex == data.EndVertex {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Error:   "Start and end vertices cannot be the same",
			Message: "Start vertex and end vertex cannot be the same according to game rules.",
		})
		return
	}
	// check if the number of edges is valid
	if data.NumberOfVertices > 1 && len(data.Edges) == 0 {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Error:   "No edges provided",
			Message: "Graph has multiple vertices but no edges, making pathfinding impossible.",
		})
		return
	}

	// Validate all edges before starting the algorithm
	for _, edge := range data.Edges {
		if edge.Source < 0 || edge.Source >= data.NumberOfVertices ||
			edge.Target < 0 || edge.Target >= data.NumberOfVertices {
			c.JSON(http.StatusBadRequest, structs.ErrorResponse{
				Error:   "Invalid edge",
				Message: fmt.Sprintf("Edge from %v to %v with a weight of %v contains invalid vertex index", edge.Source, edge.Target, edge.Weight),
			})
			return
		}
	}
	// Run Bellman-Ford algorithm using utility function
	result := utils.RunBellmanFord(data.Edges, data.NumberOfVertices, data.StartVertex)

	// Check for negative cycles
	if result.HasNegativeCycle {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Error:   "Negative cycle detected",
			Message: "The graph contains a negative cycle, making pathfinding impossible.",
		})
		return
	}

	// Check if the end vertex is reachable
	if result.Distances[data.EndVertex] == math.MaxInt32 {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Error:   "No path found",
			Message: "End vertex is not reachable from start vertex",
		})
		return
	}

	// Reconstruct the path
	path := utils.BuildPath(data.StartVertex, data.EndVertex, result.Predecessors)

	if len(path) == 0 {
		// This indicates an inconsistency: the end vertex was marked reachable,
		// but path reconstruction failed.
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Error:   "Path reconstruction error",
			Message: "Failed to reconstruct path. The end vertex might be disconnected in the predecessor tree.",
		})
		return
	}

	c.JSON(http.StatusOK, BellmanFordResponse{
		Distance: result.Distances[data.EndVertex],
		Path:     path,
	})
}
