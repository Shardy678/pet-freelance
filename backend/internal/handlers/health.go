package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func HelloWorld(c *gin.Context) {
	log.Println("wasssup")
	c.JSON(http.StatusOK, gin.H{
		"hello": "world",
	})
}
