package log

import (
	"fmt"
	"log"
	"os"
)

// re-export
const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lmsgprefix                    // move the "prefix" from the beginning of the line to before the message
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

// re-export
type (
	Logger = log.Logger
)

// re-export
var (
	New       = log.New
	Default   = log.Default
	Output    = log.Output
	SetOutput = log.SetOutput
	Flags     = log.Flags
	SetFlags  = log.SetFlags
	Prefix    = log.Prefix
	SetPrefix = log.SetPrefix
	Writer    = log.Writer
)

// re-export
var std = Default()

// re-export
// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// - rewrite these functions to use the LoggerInterface(default LevelLogger)!!!
var (
	// Print enhence std log.Printf to accept ContextLogger parameter
	// NOTE: ContextLogger parameter will be removed when print
	//
	// usage cases:
	//
	// - log.Print(log.Ctx("xxx", log.E), ...): loggerXXX.Error("", ...)
	// - log.Print(log.Ctx("xxx"), ...): loggerXXX.Info("", ...)
	// - log.Print(log.Ctx(log.E), ...): defaultLogger.Error("", ...)
	// - log.Print(...): defaultLogger.Info("", ...)
	Print = convert(false, false)
	Panic = convert(true, false)
	Fatal = convert(false, true)

	// XXXln same as `XXX`
	Println = convert(false, false)
	Panicln = convert(true, false)
	Fatalln = convert(false, true)

	// Printf enhence stdlog.Printf to accept ContextLogger parameter
	// NOTE: format contains the format directive for ContextLogger, e.g.
	//
	//     `log.Printf("%v: %d: %s", log.Ctx("xxx", log.E), code, "this is a error")`
	//
	// usage cases:
	//
	// - log.Printf(format, log.Ctx("xxx", log.E), ...): loggerXXX.Error(fmt.Sprintf(format, log.Ctx("xxx", log.E), ...))
	// - log.Printf(format, log.Ctx("xxx"), ...): loggerXXX.Info(fmt.Sprintf(format, log.Ctx("xxx"), ...))
	// - log.Printf(format, log.Ctx(log.E), ...): defaultLogger.Error(fmt.Sprintf(format, log.Ctx(log.E), ...))
	// - log.Printf(...): defaultLogger.Info(fmt.Sprintf(format, ...))
	Printf = convertf(false, false)
	Panicf = convertf(true, false)
	Fatalf = convertf(false, true)
)

func convert(isPanic, isFatal bool) func(v ...any) {
	return func(v ...any) {
		logger := Ctx()
		switch len(v) {
		case 0, 1: // nothing to do
		default:
			if ctx, ok := v[0].(ContextLogger); ok {
				logger = ctx
				v = v[1:]
			}
		}
		if isPanic || isFatal {
			logger.Level = ErrorLevel
		}
		Log(logger, logger.Level, "", v...)
		if isPanic {
			panic(fmt.Sprintln(v...))
		}
		if isFatal {
			os.Exit(1)
		}
	}
}

func convertf(isPanic, isFatal bool) func(format string, v ...any) {
	return func(format string, v ...any) {
		logger := Ctx()
		switch len(v) {
		case 0, 1: // nothing to do
		default:
			if ctx, ok := v[0].(ContextLogger); ok {
				logger = ctx
			}
		}
		if isPanic || isFatal {
			logger.Level = ErrorLevel
		}
		Log(logger, logger.Level, fmt.Sprintf(format, v...))
		if isPanic {
			panic(fmt.Sprintln(v...))
		}
		if isFatal {
			os.Exit(1)
		}
	}
}

// LevelLogger is a logger implementation, which sends al logs to provided standard output.
type LevelLogger struct {
	*Logger
	level Level
}

type LevelLoggerOption func(*LevelLogger)

func WithLogger(std *log.Logger) LevelLoggerOption {
	return func(logger *LevelLogger) {
		logger.Logger = std
	}
}

func WithLevel(level Level) LevelLoggerOption {
	return func(logger *LevelLogger) {
		logger.level = level
	}
}

// NewLevelLogger creates LevelLogger which sends all logs to specified stdlog logger.
func NewLevelLogger(opts ...LevelLoggerOption) LoggerInterface {
	levelLogger := &LevelLogger{
		Logger: Default(),
		level:  InfoLevel,
	}
	for _, opt := range opts {
		opt(levelLogger)
	}
	return levelLogger
}

func (logger *LevelLogger) Debug(msg string, keysAndValues ...interface{}) {
	if logger.level <= DebugLevel {
		var msgPart string
		if msg != "" {
			msgPart = fmt.Sprintf("msg: %s, ", msg)
		}
		logger.Print("level: debug, "+msgPart, fmt.Sprintln(keysAndValues...))
	}
}

func (logger *LevelLogger) Info(msg string, keysAndValues ...interface{}) {
	if logger.level <= InfoLevel {
		var msgPart string
		if msg != "" {
			msgPart = fmt.Sprintf("msg: %s, ", msg)
		}
		logger.Print("level: info, "+msgPart, fmt.Sprintln(keysAndValues...))
	}
}

func (logger *LevelLogger) Error(msg string, keysAndValues ...interface{}) {
	if logger.level <= ErrorLevel {
		var msgPart string
		if msg != "" {
			msgPart = fmt.Sprintf("msg: %s, ", msg)
		}
		logger.Print("level: error, "+msgPart, fmt.Sprintln(keysAndValues...))
	}
}

func Debug(msg string, keysAndValues ...interface{}) {
	defaultLogger.Debug(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...interface{}) {
	defaultLogger.Info(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	defaultLogger.Error(msg, keysAndValues...)
}

// Interface guard
var (
	_ LoggerInterface = (*LevelLogger)(nil)
)
