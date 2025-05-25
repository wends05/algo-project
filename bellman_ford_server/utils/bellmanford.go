package utils

import (
	"bellman_ford_server/structs"
	"math"
)

// HasNegativeCycle uses Bellman-Ford algorithm to detect negative cycles
// This is a utility function that can be used by both graph generation and pathfinding
func HasNegativeCycle(edges []structs.Edge, vertices int) bool {
	// Handle edge case
	if vertices <= 0 || len(edges) == 0 {
		return false
	}

	// Initialize distances
	dist := make([]int, vertices)
	for i := range dist {
		dist[i] = math.MaxInt32 // Use MaxInt32 as "infinity"
	}
	dist[0] = 0 // Start from vertex 0

	// Relax edges vertices-1 times
	for i := 0; i < vertices-1; i++ {
		relaxed := false
		for _, edge := range edges {
			if dist[edge.Source] != math.MaxInt32 && 
			   dist[edge.Source]+edge.Weight < dist[edge.Target] {
				dist[edge.Target] = dist[edge.Source] + edge.Weight
				relaxed = true
			}
		}
		// Early termination if no relaxation occurred
		if !relaxed {
			break
		}
	}

	// Check for negative cycles (one more iteration)
	for _, edge := range edges {
		if dist[edge.Source] != math.MaxInt32 && 
		   dist[edge.Source]+edge.Weight < dist[edge.Target] {
			return true // Negative cycle detected
		}
	}

	return false // No negative cycle
}

// BellmanFordResult represents the result of running Bellman-Ford algorithm
type BellmanFordResult struct {
	Distances    []int
	Predecessors []int
	HasNegativeCycle bool
}

// RunBellmanFord executes the Bellman-Ford algorithm and returns comprehensive results
// This can be used by the main Bellman-Ford endpoint for full pathfinding
func RunBellmanFord(edges []structs.Edge, vertices, startVertex int) BellmanFordResult {
	result := BellmanFordResult{
		Distances:    make([]int, vertices),
		Predecessors: make([]int, vertices),
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

// BuildPath reconstructs the shortest path from the predecessor array
func BuildPath(startVertex, endVertex int, predecessors []int) []int {
	path := []int{}
	current := endVertex

	// Traverse backwards from the end vertex using predecessors
	for current != -1 {
		path = append([]int{current}, path...) // Prepend current vertex to build path in correct order
		if current == startVertex {
			break // Stop if we reached the start vertex
		}
		current = predecessors[current]
	}

	// If the path doesn't start with the start vertex, return empty path
	if len(path) > 0 && path[0] != startVertex {
		return []int{}
	}

	return path
}
