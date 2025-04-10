package logger

import (
	"github.com/faelmori/logz"
	lgz "github.com/faelmori/logz/logger"
)

type Log struct{ lgz.Logger }

func (l *Log) Trace(msg string, fields map[string]interface{}) { logz.TraceCtx(msg, fields) }

func (l *Log) Success(msg string, fields map[string]interface{}) {
	logz.SuccessCtx(msg, fields)
}

func (l *Log) Notice(msg string, fields map[string]interface{}) {
	logz.NoticeCtx(msg, fields)
}

func (l *Log) Info(msg string, fields map[string]interface{}) { logz.InfoCtx(msg, fields) }

func (l *Log) Error(msg string, fields map[string]interface{}) { logz.ErrorCtx(msg, fields) }

func (l *Log) Warn(msg string, fields map[string]interface{}) { logz.WarnCtx(msg, fields) }

func (l *Log) Debug(msg string, fields map[string]interface{}) { logz.DebugCtx(msg, fields) }

var logger lgz.LogzLogger = lgz.NewLogger("GoLife")

func SetLogger(customLogger lgz.LogzLogger) {
	if customLogger != nil {
		logger = customLogger
	}
}

func GetLogger() lgz.LogzLogger {
	if logger == nil {
		lll := logz.NewLogger("GoLife")
		return lll
	}
	return logger
}
