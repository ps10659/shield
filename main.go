package main

import (
	"github.com/gin-gonic/gin"
)

var DB = make(map[string]string)

func main() {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})




	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

