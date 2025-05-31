package main

import (
	"bellman_ford_server/internal/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()
	
	// Setup all routes
	routes.SetupRoutes(router)

	router.Run(":" + port)
}
