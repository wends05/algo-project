package services

import (
	"bellman_ford_server/internal/models"
	"math"
)

// BellmanFordService handles all Bellman-Ford algorithm operations
type BellmanFordService struct{}

// NewBellmanFordService creates a new BellmanFordService instance
func NewBellmanFordService() *BellmanFordService {
	return &BellmanFordService{}
}

// algorithmResult represents internal algorithm results
type algorithmResult struct {
	Distances        []int
	Predecessors     []int
	HasNegativeCycle bool
}

// FindShortestPath finds the shortest path between two vertices using Bellman-Ford algorithm
func (s *BellmanFordService) FindShortestPath(input models.BellmanFordInput) (*models.BellmanFordResult, error) {
	// Run the algorithm
	result := s.runBellmanFord(input.Edges, input.NumberOfVertices, input.StartVertex)

	// Check for negative cycles
	if result.HasNegativeCycle {
		return &models.BellmanFordResult{
			HasNegativeCycle: true,
		}, nil
	}

	// Check if the end vertex is reachable
	if result.Distances[input.EndVertex] == math.MaxInt32 {
		return nil, &PathNotFoundError{
			StartVertex: input.StartVertex,
			EndVertex:   input.EndVertex,
		}
	}

	// Reconstruct the path
	path := s.buildPath(input.StartVertex, input.EndVertex, result.Predecessors)
	if len(path) == 0 {
		return nil, &PathReconstructionError{}
	}

	return &models.BellmanFordResult{
		Distance:         result.Distances[input.EndVertex],
		Path:             path,
		HasNegativeCycle: false,
	}, nil
}

// runBellmanFord executes the core Bellman-Ford algorithm
func (s *BellmanFordService) runBellmanFord(edges []models.Edge, vertices, startVertex int) algorithmResult {
	result := algorithmResult{
		Distances:        make([]int, vertices),
		Predecessors:     make([]int, vertices),
		HasNegativeCycle: false,
	}

	// Initialize distances and predecessors
	for i := range vertices {
		result.Distances[i] = math.MaxInt32
		result.Predecessors[i] = -1
	}
	result.Distances[startVertex] = 0

	// Relax edges vertices-1 times
	for range vertices - 1 {
		relaxed := false
		for _, edge := range edges {
			if result.Distances[edge.Source] != math.MaxInt32 &&
				result.Distances[edge.Source]+edge.Weight < result.Distances[edge.Target] {
				result.Distances[edge.Target] = result.Distances[edge.Source] + edge.Weight
				result.Predecessors[edge.Target] = edge.Source
				relaxed = true
			}
		}
		// Early termination if no relaxation occurred
		if !relaxed {
			break
		}
	}

	// Check for negative cycles
	for _, edge := range edges {
		if result.Distances[edge.Source] != math.MaxInt32 &&
			result.Distances[edge.Source]+edge.Weight < result.Distances[edge.Target] {
			result.HasNegativeCycle = true
			break
		}
	}

	return result
}

// buildPath reconstructs the shortest path from the predecessor array
func (s *BellmanFordService) buildPath(startVertex, endVertex int, predecessors []int) []int {
	path := []int{}
	current := endVertex

	// Traverse backwards from the end vertex using predecessors
	for current != -1 {
		path = append([]int{current}, path...) // Prepend current vertex
		if current == startVertex {
			break
		}
		current = predecessors[current]
	}

	// If the path doesn't start with the start vertex, return empty path
	if len(path) == 0 || path[0] != startVertex {
		return []int{}
	}

	return path
}
