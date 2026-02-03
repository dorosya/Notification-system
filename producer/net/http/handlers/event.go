package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func EventHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Send notification succesfully"})
}
