package controllers

import "github.com/gin-gonic/gin"

type SampleData struct {
	HelloWorld       string `json:"hello_world"`
	NumberOfVertices int    `json:"number_of_vertices"`
}

func SampleFunction(c *gin.Context) {

	data := SampleData{}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": gin.H{
			"message": "Invalid input",
			"details": err.Error(),
		}})
		return
	}

	println("Hello World:", data.HelloWorld)
	println("Number of Vertices:", data.NumberOfVertices)
	c.JSON(200, gin.H{
		"message": "Hello World",
		"data":    data,
	})
}
