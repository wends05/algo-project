package main

import (
	// gin
	c "bellman_ford_server/controllers"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := gin.Default()

	router.GET("/", c.Root)
	router.POST("/bellmanford", c.BellmanFord)

	router.GET("/graph", c.GenerateRandomGraph)

	router.Run(":" + port)
}
