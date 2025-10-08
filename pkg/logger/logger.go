package logger

import "go.uber.org/zap"

var log *zap.SugaredLogger

func InitLogger() {
	rawLogger, _ := zap.NewProduction()
	log = rawLogger.Sugar()
}

func WithFields(fields map[string]interface{}) *zap.SugaredLogger {
	return log.With(zapFields(fields)...)
}

func zapFields(fields map[string]interface{}) []interface{} {
	var z []interface{}
	for k, v := range fields {
		z = append(z, k, v)
	}
	return z
}

func Info(msg string, fields ...interface{}) {
	log.Infow(msg, fields...)
}

func Error(msg string, fields ...interface{}) {
	log.Errorw(msg, fields...)
}

func Sync() {
	log.Sync()
}

func Fatalw(msg string, fields ...interface{}) {
	log.Fatalw(msg, fields...)
}

func Default() *zap.SugaredLogger {
	return log
}
