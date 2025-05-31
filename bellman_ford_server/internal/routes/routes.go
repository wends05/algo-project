package routes

import (
	"bellman_ford_server/internal/handlers"
	"bellman_ford_server/internal/services"
	"bellman_ford_server/internal/validators"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine) {
	// Initialize services
	bellmanFordService := services.NewBellmanFordService()
	graphService := services.NewGraphService()

	// Initialize validators
	bellmanFordValidator := validators.NewBellmanFordValidator()
	graphValidator := validators.NewGraphValidator()

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	bellmanFordHandler := handlers.NewBellmanFordHandler(bellmanFordService, bellmanFordValidator)
	graphHandler := handlers.NewGraphHandler(graphService, graphValidator)

	// Health routes
	router.GET("/", healthHandler.Root)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Algorithm routes
		v1.POST("/bellmanford", bellmanFordHandler.FindShortestPath)

		// Graph generation routes
		v1.GET("/graph", graphHandler.GenerateRandomGraph)
	}

	// Backward compatibility routes (maintain existing endpoints)
	router.POST("/bellmanford", bellmanFordHandler.FindShortestPath)
	router.GET("/graph", graphHandler.GenerateRandomGraph)
}
