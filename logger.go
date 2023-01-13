package log

import (
	"sync"
)

// LoggerInterface log interface
type LoggerInterface interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// LoggerName logger name
type LoggerName string

// Level logger level
type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in production.
	DebugLevel Level = -1
	// InfoLevel is the default logging priority.
	InfoLevel Level = 0
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel Level = 2
)

var (
	// D is alias of DebugLevel
	D = DebugLevel
	// I is alias of InfoLevel
	I = InfoLevel
	// E is alias of ErrorLevel
	E = ErrorLevel
)

func (level Level) String() string {
	levelNamesLock.RLock()
	defer levelNamesLock.RUnlock()
	return levelNames[level]
}

// RegisterLevelName register level name
// TBD: support extend level?
func RegisterLevelName(level Level, name string) {
	levelNamesLock.Lock()
	defer levelNamesLock.Unlock()
	levelNames[level] = name
}

// Log utility function to log by level
func Log(logger LoggerInterface, level Level, msg string, keysAndValues ...interface{}) {
	switch level {
	case ErrorLevel:
		logger.Error(msg, keysAndValues...)
	case DebugLevel:
		logger.Debug(msg, keysAndValues...)
	default:
		logger.Info(msg, keysAndValues...)
	}
}

// GetLogger get logger from internal registry
func GetLogger(loggerName LoggerName) LoggerInterface {
	loggersLock.RLock()
	defer loggersLock.RUnlock()
	return loggers[loggerName]
}

// SetLogger set logger into internal registry
func SetLogger(loggerName LoggerName, logger LoggerInterface) {
	loggersLock.Lock()
	defer loggersLock.Unlock()
	loggers[loggerName] = logger
}

var (
	levelNames = map[Level]string{
		DebugLevel: "debug",
		InfoLevel:  "info",
		ErrorLevel: "error",
	}
	levelNamesLock sync.RWMutex
	loggers        = map[LoggerName]LoggerInterface{
		LoggerName(""):    defaultLogger,
		LoggerName("nop"): &NopLogger{},
	}
	loggersLock   sync.RWMutex
	defaultLogger = NewLevelLogger()
)
