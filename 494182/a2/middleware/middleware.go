package main

import (
	"errors"
	_"fmt"
	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware for handling errors.
func ErrorHandler(c *gin.Context) {
	defer c.Next()

	if err := c.Err; err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": err.Error(),
		})
	}
}

func main() {
	r := gin.Default()

	// Register the error handler middleware
	r.Use(ErrorHandler)

	// Sample route that simulates an error
	r.GET("/sample", func(c *gin.Context) {
		err := errors.New("sample error")
		c.Error(err)
	})

	r.Run(":8080")
}