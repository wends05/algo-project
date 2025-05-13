package structs

type BellmanFordData struct {
	NumberOfVertices int    `json:"number_of_vertices"`
	Edges            []Edge `json:"graph"`
	StartVertex      int    `json:"start_vertex"`
	EndVertex        int    `json:"end_vertex"`
	Distance         []int  `json:"distance"`
}
