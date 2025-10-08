package router

import "github.com/gin-gonic/gin"

func SetupRouter(e *gin.Engine) *gin.Engine {
	e.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return e
}
