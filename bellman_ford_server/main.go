package main

import (
	// gin
	c "bellman_ford_server/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/bellmanford", c.BellmanFord)

	router.GET("/graph", c.GenerateRandomGraph)

	router.Run(":8080")
}
