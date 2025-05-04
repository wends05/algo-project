package controllers

import (
	"github.com/gin-gonic/gin"
)

// BellmanFord input
type BellmanFordInput struct {
	Graph       [][]int `json:"graph" binding:"required"`
	StartVertex int     `json:"start_vertex" binding:"required"`
	EndVertex   int     `json:"end_vertex" binding:"required"`
}

func BellmanFord(c *gin.Context) {
	data := BellmanFordInput{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// initialize distance array

	

}
