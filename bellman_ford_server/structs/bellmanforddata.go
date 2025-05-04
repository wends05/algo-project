package structs

type BellmanFordData struct {
	Graph       [][]int `json:"graph"`
	StartVertex int     `json:"start_vertex"`
	EndVertex   int     `json:"end_vertex"`
	Distance    []int   `json:"distance"`
}
