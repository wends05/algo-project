package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"bellman_ford_server/structs"
	"bellman_ford_server/utils"

	"github.com/gin-gonic/gin"
)

type GraphResponse struct {
	NumberOfVertices int            `json:"number_of_vertices"`
	Edges            []structs.Edge `json:"edges" binding:"dive"`
	NumberOfEdges    int            `json:"number_of_edges"`
}

// GenerateRandomGraph generates a random graph with no negative cycles.
//
// parameters:
// - vertices: number of vertices in the graph (required)
// - weight_range: range of weights for the edges (optional, default is 100)
// - start_vertex: starting vertex for pathfinding (optional, default is 0)
// - end_vertex: ending vertex for pathfinding (optional, default is vertices-1)
func GenerateRandomGraph(c *gin.Context) {
	// Get number of vertices from query parameter
	verticesStr := c.Query("number_of_vertices")
	if verticesStr == "" {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Error: "number_of_vertices parameter is required",
		})
		return
	}

	// convert vertices to integer
	vertices, err := strconv.Atoi(verticesStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Error: "vertices must be a valid integer",
		})
		return
	}

	// vertices should be at least 3
	if vertices < 2 {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Error:   "vertices must be at least 3",
			Message: "Number of vertices must be at least 3 for pathfinding between distinct points.",
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

	// Get weight range from query parameter, default to 100
	weightRangeStr := c.Query("weight_range")
	if weightRangeStr == "" {
		weightRangeStr = "100"
	}

	// convert weight range to integer
	weightRange, err := strconv.Atoi(weightRangeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Error: "weight_range must be a valid integer",
		})
		return
	}

	startVertexStr := c.Query("start_vertex")
	endVertexStr := c.Query("end_vertex")


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
		NumberOfVertices: vertices,
		Edges:            edges,
		NumberOfEdges:    len(edges),
	}

	c.JSON(http.StatusOK, response)
}

// generateGraph creates a random graph with no negative cycles where all vertices are reachable
func generateGraph(vertices, weightRange int) []structs.Edge {
	var edges []structs.Edge
	edgeSet := make(map[string]bool)

	// First, create a connected backbone to ensure all vertices are reachable
	// Create a path: 0 → 1 → 2 → ... → (vertices-1)
	for i := range vertices - 1 {
		source := i
		target := i + 1
		edgeKey := strconv.Itoa(source) + "-" + strconv.Itoa(target)

		if !edgeSet[edgeKey] {
			weight := rand.Intn(2*weightRange+1) - weightRange
			edge := structs.Edge{
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
		numAdditionalEdgesToAdd := rand.Intn(maxAdditionalEdges/2 + 1)

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
			if attempts > maxAttempts {
				continue
			}

			// Generate weight and check for negative cycles
			weight := rand.Intn(2*weightRange+1) - weightRange

			// Create temporary edge to test
			testEdge := structs.Edge{
				Source: source,
				Target: target,
				Weight: weight,
			}

			// Test if adding this edge creates a negative cycle
			tempEdges := append(edges, testEdge)
			if !utils.HasNegativeCycle(tempEdges, vertices) {
				edges = tempEdges
				edgeSet[edgeKey] = true
			}
			// If it creates a negative cycle, skip this edge and try another
		}
	}

	return edges
}
