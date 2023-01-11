package log

import (
	"fmt"
	"log"
	"sync"
)

var (
	// Print
	//
	// usage cases:
	//
	// - log.Print(log.LoggerName("xxx"), log.ErrorLevel, ...): loggerXXX.Error("", ...)
	// - log.Print(log.LoggerName("xxx"), ...): loggerXXX.Info("", ...)
	// - log.Print(log.ErrorLevel, ...): defaultLogger.Error("", ...)
	// - log.Print(...): defaultLogger.Info("", ...)
	Print = convert(log.Print)
	// Println same as `Print`
	Println = convert(log.Println)

	// Printf
	//
	// usage cases:
	//
	// - log.Printf(format, log.LoggerName("xxx"), log.ErrorLevel, ...): loggerXXX.Error(fmt.Sprintf(format, ...))
	// - log.Printf(format, log.LoggerName("xxx"), ...): loggerXXX.Info(fmt.Sprintf(format, ...))
	// - log.Printf(format, log.ErrorLevel, ...): defaultLogger.Error(fmt.Sprintf(format, ...))
	// - log.Printf(...): defaultLogger.Info(fmt.Sprintf(format, ...))
	Printf = convertf(log.Printf)
)

func convert(func(v ...any)) func(v ...any) {
	return func(v ...any) {
		logger := defaultLogger
		level := InfoLevel
		switch len(v) {
		case 0, 1: // nothing to do
		default:
			if loggerName, ok := v[0].(LoggerName); ok {
				tmp := GetLogger(loggerName)
				if tmp != nil {
					logger = tmp
				} else {
					logger.Error("logger %s not found")
				}
				v = v[1:]
			}
			if l, ok := v[0].(Level); ok {
				level = l
				v = v[1:]
			}
		}
		logger.Log(level, "", v...)
	}
}

func convertf(func(format string, v ...any)) func(format string, v ...any) {
	return func(format string, v ...any) {
		logger := defaultLogger
		level := InfoLevel
		var vtmp []any
		switch len(v) {
		case 0, 1: // nothing to do
		default:
			if loggerName, ok := v[0].(LoggerName); ok {
				tmp := GetLogger(loggerName)
				if tmp != nil {
					logger = tmp
				} else {
					logger.Error("logger %s not found")
				}
				vtmp = v[1:]
			}
			if l, ok := vtmp[0].(Level); ok {
				level = l
			}
		}
		logger.Log(level, fmt.Sprintf(format, v...))
	}
}

func GetLogger(loggerName LoggerName) Logger {
	loggersLock.RLock()
	defer loggersLock.RUnlock()
	return loggers[loggerName]
}

func SetLogger(loggerName LoggerName, logger Logger) {
	loggersLock.Lock()
	defer loggersLock.Unlock()
	loggers[loggerName] = logger
}

var (
	loggers     = map[LoggerName]Logger{}
	loggersLock sync.RWMutex
)
