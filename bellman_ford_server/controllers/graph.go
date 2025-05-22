package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GraphInput struct {
	NumberOfVertices int `json:"number_of_vertices"`
}

func Graph(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Graph API is working lol",
	})
}
