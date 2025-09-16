package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "GET works"})
	})

	router.HEAD("/test", func(c *gin.Context) {
		c.Status(200)
	})

	router.Run(":8080")
}
