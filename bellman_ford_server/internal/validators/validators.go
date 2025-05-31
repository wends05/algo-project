package validators

import (
	"bellman_ford_server/internal/models"
	"fmt"
)

// BellmanFordValidator validates Bellman-Ford input
type BellmanFordValidator struct{}

// NewBellmanFordValidator creates a new validator instance
func NewBellmanFordValidator() *BellmanFordValidator {
	return &BellmanFordValidator{}
}

// ValidateInput validates the Bellman-Ford input data
func (v *BellmanFordValidator) ValidateInput(input models.BellmanFordInput) error {
	// Check if the number of vertices is valid
	if input.NumberOfVertices < 2 {
		return fmt.Errorf("number of vertices must be at least 2 for pathfinding between distinct points")
	}

	// Check if the start vertex is valid
	if input.StartVertex < 0 || input.StartVertex >= input.NumberOfVertices {
		return fmt.Errorf("invalid start vertex: %d", input.StartVertex)
	}

	// Check if the end vertex is valid
	if input.EndVertex < 0 || input.EndVertex >= input.NumberOfVertices {
		return fmt.Errorf("invalid end vertex: %d", input.EndVertex)
	}

	// Check if the start and end vertices are the same
	if input.StartVertex == input.EndVertex {
		return fmt.Errorf("start and end vertices cannot be the same")
	}

	// Check if edges are provided when needed
	if input.NumberOfVertices > 1 && len(input.Edges) == 0 {
		return fmt.Errorf("graph has multiple vertices but no edges, making pathfinding impossible")
	}

	// Validate all edges
	for _, edge := range input.Edges {
		if edge.Source < 0 || edge.Source >= input.NumberOfVertices ||
			edge.Target < 0 || edge.Target >= input.NumberOfVertices {
			return fmt.Errorf("edge from %d to %d with weight %d contains invalid vertex index", 
				edge.Source, edge.Target, edge.Weight)
		}
	}

	return nil
}

// GraphValidator validates graph generation parameters
type GraphValidator struct{}

// NewGraphValidator creates a new graph validator instance
func NewGraphValidator() *GraphValidator {
	return &GraphValidator{}
}

// ValidateGenerationParams validates graph generation parameters
func (v *GraphValidator) ValidateGenerationParams(vertices, weightRange, startVertex, endVertex int) error {
	if vertices < 3 {
		return fmt.Errorf("vertices must be at least 3")
	}

	if vertices > 100 {
		return fmt.Errorf("vertices cannot exceed 100 (performance limit)")
	}

	if startVertex < 0 || startVertex >= vertices {
		return fmt.Errorf("start_vertex %d is out of bounds [0, %d)", startVertex, vertices)
	}

	if endVertex < 0 || endVertex >= vertices {
		return fmt.Errorf("end_vertex %d is out of bounds [0, %d)", endVertex, vertices)
	}

	if startVertex == endVertex {
		return fmt.Errorf("start_vertex and end_vertex cannot be the same")
	}

	return nil
}
