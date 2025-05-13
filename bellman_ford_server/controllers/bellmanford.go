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
		c.JSON(400, gin.H{"error": gin.H{
			"message": "Invalid input",
			"details": err.Error(),
		}})
		return
	}

	// initialize distances array

	distances := make([]int, data.NumberOfVertices)
	predecessor := make([]int, data.NumberOfVertices)

	distances[data.StartVertex] = 0
	for i := range data.NumberOfVertices {
		distances[i] = math.MaxInt32 // Initialize all distances to a large number
		predecessor[i] = -1          // Initialize predecessors to -1
	}

	// distances should be [0, max, max, max, ...]
	// predecessors should be [-1, -1, -1, -1, ...]
	// set the distances of the start vertex to 0

	// for _, e := range data.Edges {
	// 	fmt.Printf("Edge: s:%v t:%v w%v\n", e.Source, e.Target, e.Weight)
	// }

	distances[data.StartVertex] = 0

	// Log initial state
	fmt.Printf("Initial distances: %v\n", distances)
	fmt.Printf("Initial predecessors: %v\n", predecessor)

	// relax edges V - 1 times
	for range data.NumberOfVertices - 1 {
		for _, edge := range data.Edges {
			fmt.Printf("Relaxing edge: %v\n", edge)
			if edge.Source < 0 || edge.Source >= data.NumberOfVertices || edge.Target < 0 || edge.Target >= data.NumberOfVertices {
				fmt.Printf("Invalid edge: %v\n", edge)

				c.JSON(400, gin.H{
					"error": "Invalid edge",
					"message": fmt.Sprintf(
						"Edge with source %v, target %v, weight %v is invalid",
						edge.Source,
						edge.Target,
						edge.Weight,
					)})
				return
			}
			if canRelax(edge, distances) {
				distances[edge.Target] = distances[edge.Source] + edge.Weight
				predecessor[edge.Target] = edge.Source
				println("Relaxing edge:", edge.Source, "->", edge.Target, "with weight", edge.Weight)
				println("New distances to vertex", edge.Target, ":", distances[edge.Target])
			} else {
				println("cannot relax this edge")
			}
		}
	}

	fmt.Printf("Distances after relaxation: %v\n", distances)
	fmt.Printf("Predecessors after relaxation: %v\n", predecessor)

	// assuming everything is already okay, check for negative cycles

	for _, edge := range data.Edges {
		if canRelax(edge, distances) {
			c.JSON(400, gin.H{"error": "Negative cycle detected"})
			return
		}
	}

	// check if the end vertex is reachable
	if distances[data.EndVertex] == math.MaxInt32 {
		c.JSON(400, gin.H{"error": "No path found", "message": "End vertex is not reachable from start vertex"})
		return
	}

	// reconstruct the path
	path := []int{}
	for current := data.EndVertex; current != -1; current = predecessor[current] {
		path = append([]int{current}, path...)
		if current == data.StartVertex {
			break
		}
	}

	// ensure that the path starts with the start vertex
	// if it is valid and reconstructed correctly
	if len(path) == 0 && data.StartVertex == data.EndVertex {
		path = []int{data.StartVertex}
	} else if len(path) == 0 && path[0] != data.StartVertex {
		c.JSON(400, gin.H{"error": "No path found", "message": "Path does not start with start vertex"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"distance": distances[data.EndVertex],
		"path":     path,
	})
}

func canRelax(edge structs.Edge, distances []int) bool {
	return distances[edge.Source] != math.MaxInt32 &&
		distances[edge.Source]+edge.Weight < distances[edge.Target]
}
