package handler

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/health", Health)
	return r
}

func Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
