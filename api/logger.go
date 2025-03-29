package api

import (
	"github.com/faelmori/logz"
)

type Logger interface {
	Info(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Debug(msg string, fields map[string]interface{})
}

type defaultLogger struct{}

func (l *defaultLogger) Info(msg string, fields map[string]interface{}) {
	logz.Info(msg, fields)
}

func (l *defaultLogger) Error(msg string, fields map[string]interface{}) {
	logz.Error(msg, fields)
}

func (l *defaultLogger) Warn(msg string, fields map[string]interface{}) {
	logz.Warn(msg, fields)
}

func (l *defaultLogger) Debug(msg string, fields map[string]interface{}) {
	logz.Debug(msg, fields)
}

var logger Logger = &defaultLogger{}

func SetLogger(customLogger Logger) {
	if customLogger != nil {
		logger = customLogger
	}
}

func GetLogger() Logger {
	return logger
}
