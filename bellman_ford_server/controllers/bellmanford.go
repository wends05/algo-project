package controllers

import (
	"bellman_ford_server/structs"
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
	EndVertex        int            `json:"end_vertex" binding:"required"`
}

func BellmanFord(c *gin.Context) {

	data := BellmanFordInput{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"message": "Invalid input",
			"details": err.Error(),
		}})
		return
	}

	// --- Input Validation (more comprehensive) ---

	// Check if the number of vertices is valid
	if data.NumberOfVertices < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid graph size",
			"message": "Number of vertices must be at least 2 for pathfinding between distinct points.",
		})
		return
	}
	// Check if the number of edges is valid
	if data.StartVertex < 0 || data.StartVertex >= data.NumberOfVertices {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start vertex"})
		return
	}
	// Check if the start vertex is valid
	if data.EndVertex < 0 || data.EndVertex >= data.NumberOfVertices {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end vertex"})
		return
	}
	// Check if the start and end vertices are the same
	if data.StartVertex == data.EndVertex {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid game state",
			"message": "Start vertex and end vertex cannot be the same according to game rules.",
		})
		return
	}
	// check if the number of edges is valid
	if data.NumberOfVertices > 1 && len(data.Edges) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No edges provided",
			"message": "Graph has multiple vertices but no edges, making pathfinding impossible.",
		})
		return
	}

	// Validate all edges before starting the algorithm
	for _, edge := range data.Edges {
		if edge.Source < 0 || edge.Source >= data.NumberOfVertices ||
			edge.Target < 0 || edge.Target >= data.NumberOfVertices {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid edge",
				"message": fmt.Sprintf("Edge with source %v, target %v, weight %v contains invalid vertex index", edge.Source, edge.Target, edge.Weight),
			})
			return
		}
	}

	// initialize distances array

	distances := make([]int, data.NumberOfVertices)
	predecessor := make([]int, data.NumberOfVertices)

	// distances[data.StartVertex] = 0 // This line is redundant
	for i := range data.NumberOfVertices {
		distances[i] = math.MaxInt32
		predecessor[i] = -1
	}

	// set the distance to the start vertex to 0
	distances[data.StartVertex] = 0

	// Log initial state
	fmt.Printf("Initial distances: %v\n", distances)
	fmt.Printf("Initial predecessors: %v\n", predecessor)

	// relax edges V - 1 times
	for range data.NumberOfVertices - 1 {

		relaxedInThisIteration := false

		for _, edge := range data.Edges {
			fmt.Printf("Relaxing edge: %v\n", edge)
			if distances[edge.Source] != math.MaxInt32 &&
				distances[edge.Source]+edge.Weight < distances[edge.Target] {
				distances[edge.Target] = distances[edge.Source] + edge.Weight
				predecessor[edge.Target] = edge.Source
				println("Relaxing edge:", edge.Source, "->", edge.Target, "with weight", edge.Weight)
				println("New distances to vertex", edge.Target, ":", distances[edge.Target])

				relaxedInThisIteration = true
			} else {
				println("cannot relax this edge")
			}
		}

		// If no edges were relaxed in an entire iteration, distances have converged
		if !relaxedInThisIteration {
			fmt.Println("Distances converged early.")
			break
		}
	}

	fmt.Printf("Distances after relaxation: %v\n", distances)
	fmt.Printf("Predecessors after relaxation: %v\n", predecessor)

	// assuming everything is already okay, check for negative cycles

	for _, edge := range data.Edges {
		if distances[edge.Source] != math.MaxInt32 &&
			distances[edge.Source]+edge.Weight < distances[edge.Target] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Negative cycle detected"})
			return
		}
	}

	// check if the end vertex is reachable
	if distances[data.EndVertex] == math.MaxInt32 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No path found",
			"message": "End vertex is not reachable from start vertex",
		})
		return
	}

	// reconstruct the path
	path := buildPath(data.StartVertex, data.EndVertex, predecessor)

	if len(path) == 0 {
		// This indicates an inconsistency: the end vertex was marked reachable,
		// but path reconstruction failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Path reconstruction error",
			"message": "Failed to reconstruct path. The end vertex might be disconnected in the predecessor tree.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"distance": distances[data.EndVertex],
		"path":     path,
	})
}


// buildPath reconstructs the shortest path from the predecessor array.
func buildPath(startVertex, endVertex int, predecessor []int) []int {
	path := []int{}
	current := endVertex

	// Traverse backwards from the end vertex using predecessors
	for current != -1 {
		path = append([]int{current}, path...) // Prepend current vertex to build path in correct order
		if current == startVertex {
			break // Stop if we reached the start vertex
		}
		current = predecessor[current]
	}

	// If the path doesn't start with the start vertex (and it should if reachable),
	// it indicates an issue or an unreachable path segment.
	// This check is mostly for robustness; the algorithm should prevent this.
	if len(path) > 0 && path[0] != startVertex {
		return []int{} // Return empty path if it's malformed
	}

	return path
}
