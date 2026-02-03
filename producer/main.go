package main

import (
	"net/http"
	handlers "notification-system/producer/net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})
	router.POST("/event", handlers.EventHandler)
	router.Run(":8080")
}
