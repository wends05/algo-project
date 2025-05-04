package main

import (
	// gin
	"bellman_ford_server/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/bellmanford", controllers.BellmanFord)
	router.POST("/bellmanford", controllers.BellmanFord)

	router.Run(":8080")
}
