package services

import "fmt"

// PathNotFoundError represents an error when no path exists between vertices
type PathNotFoundError struct {
	StartVertex int
	EndVertex   int
}

func (e *PathNotFoundError) Error() string {
	return fmt.Sprintf("no path found from vertex %d to vertex %d", e.StartVertex, e.EndVertex)
}

// PathReconstructionError represents an error during path reconstruction
type PathReconstructionError struct{}

func (e *PathReconstructionError) Error() string {
	return "failed to reconstruct path - the end vertex might be disconnected in the predecessor tree"
}

// InvalidGraphError represents an error with invalid graph parameters
type InvalidGraphError struct {
	Message string
}

func (e *InvalidGraphError) Error() string {
	return e.Message
}
