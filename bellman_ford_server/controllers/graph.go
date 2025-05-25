package controllers

import (
	"bellman_ford_server/structs"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Graph struct {
	NumberOfVertices int            `json:"number_of_vertices"`
	Edges            []structs.Edge `json:"edges" binding:"dive"`
	NumberOfEdges    int            `json:"number_of_edges"`
}

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

	edges := generateGraph(vertices, weightRange, startVertex, endVertex)

	c.JSON(http.StatusOK, Graph{
		NumberOfVertices: vertices,
		Edges:            edges,
		NumberOfEdges:    len(edges),
	})
}

func generateGraph(vertices, weightRange, startVertex, endVertex int) []structs.Edge {

	// With max 3 outgoing edges per vertex, maximum possible edges is 3 * vertices
	maxPossibleEdges := min(vertices*3, vertices*(vertices-1)) // Respect both constraints
	targetNumEdges := rand.Intn(maxPossibleEdges/2 + 1) + maxPossibleEdges/4 // Target 25-75% of max possible

	edges := make([]structs.Edge, 0, targetNumEdges)
	edgeSet := make(map[string]bool)

	// Create levels ensuring start vertex has lower level than end vertex
	levels := make([]int, vertices)
	for i := range vertices {
		levels[i] = i
	}

	// Ensure start vertex has a lower level than end vertex for guaranteed path
	if levels[startVertex] >= levels[endVertex] {
		levels[startVertex], levels[endVertex] = levels[endVertex], levels[startVertex]
	}

	// Create a guaranteed path from start to end vertex
	pathVertices := []int{startVertex}

	// Find vertices that can be part of the path (have levels between start and end)
	var intermediateVertices []int
	for i := range vertices {
		if i != startVertex && i != endVertex &&
			levels[i] > levels[startVertex] && levels[i] < levels[endVertex] {
			intermediateVertices = append(intermediateVertices, i)
		}
	}

	// Ensure at least one intermediate vertex is REQUIRED (no direct start->end edge allowed)
	numIntermediateVertices := 0
	if len(intermediateVertices) > 0 {
		// Always select at least 1 intermediate vertex to prevent direct start->end edge
		minVertices := 1
		maxVertices := min(len(intermediateVertices), 3) // Limit to 3 intermediate vertices
		numIntermediateVertices = rand.Intn(maxVertices-minVertices+1) + minVertices // 1 to maxVertices
		
		// Shuffle intermediate vertices to randomize selection
		rand.Shuffle(len(intermediateVertices), func(i, j int) {
			intermediateVertices[i], intermediateVertices[j] = intermediateVertices[j], intermediateVertices[i]
		})

		// Add selected intermediate vertices to path
		for i := range numIntermediateVertices {
			pathVertices = append(pathVertices, intermediateVertices[i])
		}
	} else {
		// If no intermediate vertices available, create one by adjusting levels
		// Find a vertex that can serve as intermediate
		var candidateVertex int = -1
		for i := range vertices {
			if i != startVertex && i != endVertex {
				candidateVertex = i
				break
			}
		}
		
		if candidateVertex != -1 {
			// Adjust the candidate vertex level to be between start and end
			levels[candidateVertex] = (levels[startVertex] + levels[endVertex]) / 2
			pathVertices = append(pathVertices, candidateVertex)
		}
	}

	// Add end vertex to complete the path
	pathVertices = append(pathVertices, endVertex)

	// Sort path vertices by their levels to ensure correct order
	sort.Slice(pathVertices, func(i, j int) bool {
		return levels[pathVertices[i]] < levels[pathVertices[j]]
	})

	// Create edges for the guaranteed path
	for i := range len(pathVertices)-1 {
		source := pathVertices[i]
		target := pathVertices[i+1]
		weight := rand.Intn(2*weightRange+1) - weightRange

		edgeKey := fmt.Sprintf("%d-%d", source, target)
		if !edgeSet[edgeKey] {
			newEdge := structs.Edge{
				Source: source,
				Target: target,
				Weight: weight,
			}
			edges = append(edges, newEdge)
			edgeSet[edgeKey] = true
		}
	}

	// Track outgoing edge count for each vertex (max 3 per vertex)
	outgoingEdgeCount := make([]int, vertices)
	
	// Count edges from the guaranteed path
	for _, edge := range edges {
		outgoingEdgeCount[edge.Source]++
	}

	// Generate additional random edges to reach target number
	maxAttempts := max(targetNumEdges*vertices*2, 100)
	attempts := 0

	for len(edges) < targetNumEdges && attempts < maxAttempts {
		attempts++

		u := rand.Intn(vertices)
		v := rand.Intn(vertices)

		if u == v { // An edge cannot connect a node to itself.
			continue
		}

		var source, destination int

		if levels[u] < levels[v] {
			source = u
			destination = v
		} else {
			source = v
			destination = u
		}

		// Check if source vertex already has 3 outgoing edges
		if outgoingEdgeCount[source] >= 3 {
			continue // Skip this edge, source vertex is at max capacity
		}

		// Prevent direct edge from start vertex to end vertex
		if source == startVertex && destination == endVertex {
			continue // Skip this edge, no direct start->end allowed
		}

		edgeKey := fmt.Sprintf("%d-%d", source, destination)
		if edgeSet[edgeKey] {
			continue // This edge already exists, try picking new nodes.
		}

		weight := rand.Intn(2*weightRange+1) - weightRange

		newEdge := structs.Edge{
			Source: source,
			Target: destination,
			Weight: weight,
		}

		edges = append(edges, newEdge)
		edgeSet[edgeKey] = true // Mark this edge as added.
		outgoingEdgeCount[source]++ // Increment outgoing edge count for source vertex
	}

	return edges
}

// Helper function to find minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
