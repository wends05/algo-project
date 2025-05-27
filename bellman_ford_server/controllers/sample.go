package controllers

import "github.com/gin-gonic/gin"

func Root(c* gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to the Bellman-Ford Server! Use /bellmanford for pathfinding or /graph for random graph generation.",
	})
}
