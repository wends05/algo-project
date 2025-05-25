package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"bellman_ford_server/structs"

	"github.com/gin-gonic/gin"
)

type Edge struct {
	Source int `json:"source"`
	Target int `json:"target"`
	Weight int `json:"weight"`
}

type GraphResponse struct {
	Vertices      int    `json:"vertices"`
	Edges         []Edge `json:"edges"`
	NumberOfEdges int    `json:"number_of_edges"`
}

// GenerateRandomGraph generates a random graph with no negative cycles

// parameters:
// - vertices: number of vertices in the graph (required)
// - weight_range: range of weights for the edges (optional, default is 100)
// - start_vertex: starting vertex for pathfinding (optional, default is 0)
// - end_vertex: ending vertex for pathfinding (optional, default is vertices-1)

func GenerateRandomGraph(c *gin.Context) {
	// Get number of vertices from query parameter
	verticesStr := c.Query("vertices")
	if verticesStr == "" {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Error: "vertices parameter is required",
		})
		return
	}

	// Get weight range from query parameter, default to 100
	weightRangeStr := c.Query("weight_range")
	if weightRangeStr == "" {
		weightRangeStr = "100"
	}

	startVertexStr := c.Query("start_vertex")
	endVertexStr := c.Query("end_vertex")

	// convert vertices to integer
	vertices, err := strconv.Atoi(verticesStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Error: "vertices must be a valid integer",
		})
		return
	}

	// convert weight range to integer
	weightRange, err := strconv.Atoi(weightRangeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Error: "weight_range must be a valid integer",
		})
		return
	}

	// vertices should be at least 3
	if vertices < 2 {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Error:   "vertices must be at least 1",
			Message: "Number of vertices must be at least 2 for pathfinding between distinct points.",
		})
		return
	}

	// for simplicity limit the number of vertices to 100
	if vertices > 100 {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Error: "vertices cannot exceed 100 (performance limit)",
		})
		return
	}

	// Parse start and end vertices, if provided
	var startVertex, endVertex int
	if startVertexStr != "" {
		startVertex, err = strconv.Atoi(startVertexStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, structs.ErrorResponse{
				Error: "start_vertex must be a valid integer",
			})
			return
		}
		if startVertex < 0 || startVertex >= vertices {
			c.JSON(http.StatusBadRequest, structs.ErrorResponse{
				Error:   "Invalid start_vertex",
				Message: fmt.Sprintf("start_vertex %d is out of bounds [0, %d)", startVertex, vertices),
			})
			return
		}
	} else {
		// Default to 0 if not provided
		startVertex = 0
	}

	if endVertexStr != "" {
		endVertex, err = strconv.Atoi(endVertexStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, structs.ErrorResponse{
				Error: "end_vertex must be a valid integer",
			})
			return
		}
		if endVertex < 0 || endVertex >= vertices {
			c.JSON(http.StatusBadRequest, structs.ErrorResponse{
				Error:   "Invalid end_vertex",
				Message: fmt.Sprintf("end_vertex %d is out of bounds [0, %d)", endVertex, vertices),
			})
			return
		}
	} else {
		// Default to last vertex if not provided
		endVertex = vertices - 1
	}

	if startVertex == endVertex {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Error:   "start_vertex and end_vertex cannot be the same",
			Message: "Please provide distinct start and end vertices for pathfinding.",
		})
		return
	}

	// log all parameters
	fmt.Printf("Generating graph with %d vertices, weight range %d, start vertex %d, end vertex %d\n",
		vertices, weightRange, startVertex, endVertex)

	// Generate random graph
	edges := generateGraph(vertices, weightRange)

	response := GraphResponse{
		Vertices:      vertices,
		Edges:         edges,
		NumberOfEdges: len(edges),
	}

	c.JSON(http.StatusOK, response)
}

// generateGraph creates a random graph with no negative cycles where all vertices are reachable
func generateGraph(vertices, weightRange int) []Edge {
	var edges []Edge
	edgeSet := make(map[string]bool)

	// First, create a connected backbone to ensure all vertices are reachable
	// Create a path: 0 → 1 → 2 → ... → (vertices-1)
	for i := range vertices - 1 {
		source := i
		target := i + 1
		edgeKey := strconv.Itoa(source) + "-" + strconv.Itoa(target)

		if !edgeSet[edgeKey] {
			weight := rand.Intn(2*weightRange+1) - weightRange // Weight between -100 and 100
			edge := Edge{
				Source: source,
				Target: target,
				Weight: weight,
			}
			edges = append(edges, edge)
			edgeSet[edgeKey] = true
		}
	}

	// Calculate additional edges to add (between 0 and remaining possible edges)

	maxPossibleEdges := vertices * (vertices - 1)
	currentEdges := len(edges)
	maxAdditionalEdges := maxPossibleEdges - currentEdges

	if maxAdditionalEdges > 0 {
		numAdditionalEdgesToAdd := rand.Intn(maxAdditionalEdges/2 + 1) // Add up to half of remaining possible edges

		// Add random additional edges
		for range numAdditionalEdgesToAdd {
			var source, target int
			var edgeKey string

			// Find a unique edge
			attempts := 0
			const maxAttempts = 100

			for {
				source = rand.Intn(vertices)
				target = rand.Intn(vertices)

				// Avoid self-loops
				if source == target {
					attempts++
					if attempts > maxAttempts {
						break
					}
					continue
				}

				// Check if the edge already exists
				edgeKey = strconv.Itoa(source) + "-" + strconv.Itoa(target)
				if !edgeSet[edgeKey] {
					break
				}

				// Avoid infinite loop if too many attempts
				attempts++
				if attempts > maxAttempts {
					break
				}
			}

			// If we couldn't find a unique edge after many attempts, skip this iteration
			if attempts > 100 {
				continue
			}

			// Generate weight that can be negative or positive
			weight := rand.Intn(2*weightRange+1) - weightRange // Weight between -100 and 100
			edge := Edge{
				Source: source,
				Target: target,
				Weight: weight,
			}

			edges = append(edges, edge)
			edgeSet[edgeKey] = true
		}
	}

	return edges
}
