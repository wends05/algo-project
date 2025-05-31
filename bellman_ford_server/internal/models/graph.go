package models

// Graph represents a graph structure with vertices and edges
type Graph struct {
	NumberOfVertices int    `json:"number_of_vertices"`
	Edges            []Edge `json:"edges" binding:"dive"`
	NumberOfEdges    int    `json:"number_of_edges"`
}

// Edge represents a directed edge in the graph
type Edge struct {
	Source int `json:"source"`
	Target int `json:"target"`
	Weight int `json:"weight"`
}

// BellmanFordInput represents the input data needed to run the Bellman-Ford algorithm
type BellmanFordInput struct {
	NumberOfVertices int    `json:"number_of_vertices"`
	Edges            []Edge `json:"edges" binding:"dive"`
	StartVertex      int    `json:"start_vertex" binding:"min=0"`
	EndVertex        int    `json:"end_vertex" binding:"min=0"`
}

// BellmanFordResult represents the result of running Bellman-Ford algorithm
type BellmanFordResult struct {
	Distance         int   `json:"distance"`
	Path             []int `json:"path"`
	HasNegativeCycle bool  `json:"has_negative_cycle"`
}

// GraphGenerationParams represents parameters for graph generation
type GraphGenerationParams struct {
	Vertices    int
	WeightRange int
	StartVertex int
	EndVertex   int
}
