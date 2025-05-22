package main

import (
	// gin
	c "bellman_ford_server/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/bellmanford", c.BellmanFord)
	router.POST("/sample", c.SampleFunction)

	router.GET("/graph", c.Graph)

	router.Run(":8080")
}
