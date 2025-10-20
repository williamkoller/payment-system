package logger

import (
	"log"
	_ "os"
	"sync"

	"go.uber.org/zap"
	_ "go.uber.org/zap/zapcore"
)

var (
	logger  *zap.SugaredLogger
	once    sync.Once
	initErr error
)

func InitLogger(mode string) error {
	once.Do(func() {
		var zapLogger *zap.Logger

		switch mode {
		case "dev":
			zapLogger, initErr = zap.NewDevelopment()
		default:
			cfg := zap.NewProductionConfig()
			cfg.OutputPaths = []string{"stdout"} // Logs go to stdout
			cfg.ErrorOutputPaths = []string{"stderr"}
			zapLogger, initErr = cfg.Build()
		}

		if initErr != nil {
			log.Printf("Failed to initialize zap logger: %v", initErr)
			return
		}

		logger = zapLogger.Sugar()
	})

	return initErr
}

func WithFields(fields map[string]interface{}) *zap.SugaredLogger {
	if logger == nil {
		log.Println("Logger not initialized")
		return nil
	}
	return logger.With(zapFields(fields)...)
}

func zapFields(fields map[string]interface{}) []interface{} {
	z := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		z = append(z, k, v)
	}
	return z
}

func Info(msg string, fields ...interface{}) {
	if logger != nil {
		logger.Infow(msg, fields...)
	}
}

func Error(msg string, fields ...interface{}) {
	if logger != nil {
		logger.Errorw(msg, fields...)
	}
}

func Fatal(msg string, fields ...interface{}) {
	if logger != nil {
		logger.Fatalw(msg, fields...)
	}
}

func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}

func Default() *zap.SugaredLogger {
	if logger == nil {
		log.Println("Logger not initialized")
		return nil
	}
	return logger
}
