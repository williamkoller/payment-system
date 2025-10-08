package middleware

import "github.com/gin-gonic/gin"

func Middlewares(e *gin.Engine) {
	e.Use(gin.Recovery())
	e.Use(ZapLoggerMiddleware())
}
