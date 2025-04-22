package logger

import (
	"fmt"
	l "github.com/faelmori/logz"
	"reflect"
	"runtime"
	"strings"
)

type gLog struct {
	l.Logger
	gLogLevel LogType
}

var (
	// debug is a boolean that indicates whether to log debug messages.
	debug bool
	// g is the global logger instance.
	g *gLog = &gLog{
		Logger:    l.GetLogger("GoLife - Test"),
		gLogLevel: LogTypeInfo,
	}
)

func init() {
	// Set the debug flag to true for testing purposes.
	debug = false
	// Initialize the global logger instance with a default logger.
	if g.Logger == nil {
		g = &gLog{
			Logger:    l.GetLogger("GoLife - Test"),
			gLogLevel: LogTypeInfo,
		}
	}
}

type LogType string

const (
	LogTypeNotice  LogType = "notice"
	LogTypeInfo    LogType = "info"
	LogTypeDebug   LogType = "debug"
	LogTypeError   LogType = "error"
	LogTypeWarn    LogType = "warn"
	LogTypeFatal   LogType = "fatal"
	LogTypePanic   LogType = "panic"
	LogTypeSuccess LogType = "success"
)

// SetDebug is a function that sets the debug flag for logging.
func SetDebug(d bool) { debug = d }

// LogObjLogger is a function that logs messages with the specified log type.
func LogObjLogger[T any](obj *T, logType string, messages ...string) {
	if obj == nil {
		g.Error(fmt.Sprintf("log object (%s) is nil", reflect.TypeFor[T]()), map[string]any{
			"context":  "Log",
			"logType":  logType,
			"object":   obj,
			"msg":      messages,
			"showData": true,
		})
		return
	}
	// Check if the object has a logger field with reflection
	objValueLogger := reflect.ValueOf(obj).Elem().FieldByName("Logger")
	if !objValueLogger.IsValid() {
		g.Error(fmt.Sprintf("log object (%s) does not have a logger field", reflect.TypeFor[T]()), map[string]any{
			"context":  "Log",
			"logType":  logType,
			"object":   obj,
			"msg":      messages,
			"showData": true,
		})
		return
	} else {
		var lgr l.Logger
		objValueLogger = objValueLogger.Convert(reflect.TypeFor[l.Logger]())
		if objValueLogger.IsNil() {
			objValueLogger = reflect.ValueOf(g.Logger)
		}
		if lgr = objValueLogger.Interface().(l.Logger); lgr == nil {
			lgr = g.Logger
		}
		pc, file, line, ok := runtime.Caller(1)
		if !ok {
			lgr.Error("Log: unable to get caller information", nil)
			return
		}
		funcName := runtime.FuncForPC(pc).Name()
		ctxMessageMap := map[string]any{
			"context":  funcName,
			"file":     file,
			"line":     line,
			"showData": debug,
		}
		fullMessage := strings.Join(messages, " ")
		logType = strings.ToLower(logType)
		if logType != "" {
			if reflect.TypeOf(logType).ConvertibleTo(reflect.TypeFor[LogType]()) {
				lType := LogType(logType)
				ctxMessageMap["logType"] = logType
				logging(lgr, lType, fullMessage, ctxMessageMap)
			} else {
				lgr.Error(fmt.Sprintf("logType (%s) is not valid", logType), ctxMessageMap)
			}
		} else {
			lgr.Info(fmt.Sprintf("%s", fullMessage), ctxMessageMap)
		}
	}
}

// Log is a function that logs messages with the specified log type and caller information.
func Log(logType string, messages ...string) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		g.Error("Log: unable to get caller information", nil)
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	ctxMessageMap := map[string]any{
		"context":  funcName,
		"file":     file,
		"line":     line,
		"showData": debug,
	}
	fullMessage := strings.Join(messages, " ")
	logType = strings.ToLower(logType)
	if logType != "" {
		if reflect.TypeOf(logType).ConvertibleTo(reflect.TypeFor[LogType]()) {
			lType := LogType(logType)
			ctxMessageMap["logType"] = logType
			logging(g.Logger, lType, fullMessage, ctxMessageMap)
		} else {
			g.Error(fmt.Sprintf("logType (%s) is not valid", logType), ctxMessageMap)
		}
	} else {
		g.Info(fmt.Sprintf("%s", fullMessage), ctxMessageMap)
	}
}

// logging is a helper function that logs messages with the specified log type.
func logging(lgr l.Logger, lType LogType, fullMessage string, ctxMessageMap map[string]interface{}) {
	debugCtx := debug
	if !debugCtx {
		if lType == "error" || lType == "fatal" || lType == "panic" || lType == "debug" {
			// If debug is false, set the debug value based on the logType
			debugCtx = true
		} else {
			debugCtx = false
		}
	}
	ctxMessageMap["showData"] = debugCtx
	switch lType {
	case LogTypeInfo:
		lgr.Info(fmt.Sprintf("%s", fullMessage), ctxMessageMap)
	case LogTypeDebug:
		lgr.Debug(fmt.Sprintf("%s", fullMessage), ctxMessageMap)
	case LogTypeError:
		lgr.Error(fmt.Sprintf("%s", fullMessage), ctxMessageMap)
	case LogTypeWarn:
		lgr.Warn(fmt.Sprintf("%s", fullMessage), ctxMessageMap)
	case LogTypeNotice:
		lgr.Notice(fmt.Sprintf("%s", fullMessage), ctxMessageMap)
	case LogTypeSuccess:
		lgr.Success(fmt.Sprintf("%s", fullMessage), ctxMessageMap)
	case LogTypeFatal:
		lgr.FatalC(fmt.Sprintf("%s", fullMessage), ctxMessageMap)
	case LogTypePanic:
		lgr.Panic(fmt.Sprintf("%s", fullMessage), ctxMessageMap)
	default:
		lgr.Info(fmt.Sprintf("%s", fullMessage), ctxMessageMap)
	}
	debugCtx = debug
}
