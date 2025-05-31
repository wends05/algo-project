package services

import (
	"bellman_ford_server/internal/models"
	"fmt"
	"math/rand"
	"sort"
)

// GraphService handles graph generation operations
type GraphService struct{}

// NewGraphService creates a new GraphService instance
func NewGraphService() *GraphService {
	return &GraphService{}
}

// GenerateRandomGraph creates a random graph with the specified parameters
func (s *GraphService) GenerateRandomGraph(params models.GraphGenerationParams) (*models.Graph, error) {
	if err := s.validateParams(params); err != nil {
		return nil, err
	}

	edges := s.generateEdges(params)
	
	return &models.Graph{
		NumberOfVertices: params.Vertices,
		Edges:            edges,
		NumberOfEdges:    len(edges),
	}, nil
}

// validateParams validates the graph generation parameters
func (s *GraphService) validateParams(params models.GraphGenerationParams) error {
	if params.Vertices < 3 {
		return &InvalidGraphError{
			Message: "number of vertices must be at least 3 for pathfinding between distinct points",
		}
	}

	if params.Vertices > 100 {
		return &InvalidGraphError{
			Message: "vertices cannot exceed 100 (performance limit)",
		}
	}

	if params.StartVertex < 0 || params.StartVertex >= params.Vertices {
		return &InvalidGraphError{
			Message: fmt.Sprintf("start_vertex %d is out of bounds [0, %d)", params.StartVertex, params.Vertices),
		}
	}

	if params.EndVertex < 0 || params.EndVertex >= params.Vertices {
		return &InvalidGraphError{
			Message: fmt.Sprintf("end_vertex %d is out of bounds [0, %d)", params.EndVertex, params.Vertices),
		}
	}

	if params.StartVertex == params.EndVertex {
		return &InvalidGraphError{
			Message: "start_vertex and end_vertex cannot be the same",
		}
	}

	return nil
}

// generateEdges creates the edges for the graph
func (s *GraphService) generateEdges(params models.GraphGenerationParams) []models.Edge {
	vertices := params.Vertices
	weightRange := params.WeightRange
	startVertex := params.StartVertex
	endVertex := params.EndVertex

	// With max 3 outgoing edges per vertex, maximum possible edges is 3 * vertices
	maxPossibleEdges := min(vertices*3, vertices*(vertices-1))
	targetNumEdges := rand.Intn(maxPossibleEdges/2+1) + maxPossibleEdges/4 // Target 25-75% of max possible

	edges := make([]models.Edge, 0, targetNumEdges)
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

	// Create guaranteed path from start to end vertex
	pathVertices := s.createGuaranteedPath(vertices, startVertex, endVertex, levels)

	// Create edges for the guaranteed path
	for i := range len(pathVertices) - 1 {
		source := pathVertices[i]
		target := pathVertices[i+1]
		weight := rand.Intn(2*weightRange+1) - weightRange

		edgeKey := fmt.Sprintf("%d-%d", source, target)
		if !edgeSet[edgeKey] {
			newEdge := models.Edge{
				Source: source,
				Target: target,
				Weight: weight,
			}
			edges = append(edges, newEdge)
			edgeSet[edgeKey] = true
		}
	}

	// Generate additional random edges
	edges = s.addRandomEdges(edges, edgeSet, levels, targetNumEdges, vertices, weightRange, startVertex, endVertex)

	return edges
}

// createGuaranteedPath creates a guaranteed path from start to end vertex
func (s *GraphService) createGuaranteedPath(vertices, startVertex, endVertex int, levels []int) []int {
	pathVertices := []int{startVertex}

	// Find vertices that can be part of the path
	var intermediateVertices []int
	for i := range vertices {
		if i != startVertex && i != endVertex &&
			levels[i] > levels[startVertex] && levels[i] < levels[endVertex] {
			intermediateVertices = append(intermediateVertices, i)
		}
	}

	// Ensure at least one intermediate vertex
	numIntermediateVertices := 0
	if len(intermediateVertices) > 0 {
		minVertices := 1
		maxVertices := min(len(intermediateVertices), 3)
		numIntermediateVertices = rand.Intn(maxVertices-minVertices+1) + minVertices

		// Shuffle and select intermediate vertices
		rand.Shuffle(len(intermediateVertices), func(i, j int) {
			intermediateVertices[i], intermediateVertices[j] = intermediateVertices[j], intermediateVertices[i]
		})

		for i := range numIntermediateVertices {
			pathVertices = append(pathVertices, intermediateVertices[i])
		}
	} else {
		// Create intermediate vertex by adjusting levels
		var candidateVertex int = -1
		for i := range vertices {
			if i != startVertex && i != endVertex {
				candidateVertex = i
				break
			}
		}

		if candidateVertex != -1 {
			levels[candidateVertex] = (levels[startVertex] + levels[endVertex]) / 2
			pathVertices = append(pathVertices, candidateVertex)
		}
	}

	// Add end vertex and sort by levels
	pathVertices = append(pathVertices, endVertex)
	sort.Slice(pathVertices, func(i, j int) bool {
		return levels[pathVertices[i]] < levels[pathVertices[j]]
	})

	return pathVertices
}

// addRandomEdges adds additional random edges to reach target number
func (s *GraphService) addRandomEdges(edges []models.Edge, edgeSet map[string]bool, levels []int, targetNumEdges, vertices, weightRange, startVertex, endVertex int) []models.Edge {
	// Track outgoing edge count for each vertex (max 3 per vertex)
	outgoingEdgeCount := make([]int, vertices)

	// Count edges from the guaranteed path
	for _, edge := range edges {
		outgoingEdgeCount[edge.Source]++
	}

	maxAttempts := max(targetNumEdges*vertices*2, 100)
	attempts := 0

	for len(edges) < targetNumEdges && attempts < maxAttempts {
		attempts++

		u := rand.Intn(vertices)
		v := rand.Intn(vertices)

		if u == v {
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

		// Check constraints
		if outgoingEdgeCount[source] >= 3 {
			continue
		}

		if source == startVertex && destination == endVertex {
			continue
		}

		edgeKey := fmt.Sprintf("%d-%d", source, destination)
		if edgeSet[edgeKey] {
			continue
		}

		weight := rand.Intn(2*weightRange+1) - weightRange

		newEdge := models.Edge{
			Source: source,
			Target: destination,
			Weight: weight,
		}

		edges = append(edges, newEdge)
		edgeSet[edgeKey] = true
		outgoingEdgeCount[source]++
	}

	return edges
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
