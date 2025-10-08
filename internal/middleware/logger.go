package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/payment-system/pkg/logger"
	"github.com/williamkoller/payment-system/pkg/ulid"
	"go.uber.org/zap"
)

const loggerKey = "logger"

func ZapLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := ulid.NewULID()
		ctxLogger := logger.WithFields(map[string]interface{}{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
		})

		c.Set(loggerKey, ctxLogger)

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		ctxLogger.Infow("request completed",
			"status", status,
			"duration", duration.String(),
		)
	}
}

func FromContext(c *gin.Context) *zap.SugaredLogger {
	if log, exists := c.Get(loggerKey); exists {
		if lgr, ok := log.(*zap.SugaredLogger); ok {
			return lgr
		}
	}
	return logger.Default()
}
