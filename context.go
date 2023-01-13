package log

import (
	"fmt"
)

// Ctx returna a conetxt logger
// usage:
// - log.Ctx("xxx", log.ErrorLevel): use xxx logger in error level
// - log.Ctx("xxx"): use xxx logger in info level
// - log.Ctx(log.DebugLevel): use default logger in debug level
func Ctx(ps ...any) ContextLogger {
	var opts []ContextLoggerOption
	for _, p := range ps {
		switch p := p.(type) {
		case string:
			opts = append(opts, WithContextLoggerName(p))
		case int8:
			opts = append(opts, WithContextLevel(p))
		case LoggerName:
			opts = append(opts, WithContextLoggerName(string(p)))
		case Level:
			opts = append(opts, WithContextLevel(int8(p)))
		case ContextLoggerOption:
			opts = append(opts, p)
		default:
			panic(fmt.Errorf("log.Ctx not support parameter type: %T", p))
		}
	}
	ctx := ContextLogger{
		Level: InfoLevel,
	}
	for _, opt := range opts {
		opt(&ctx)
	}
	return ctx
}

// ContextLoggerOption specify context logger option
type ContextLoggerOption func(*ContextLogger)

// WithContextLoggerName specify logger name
func WithContextLoggerName(name string) ContextLoggerOption {
	return func(ctx *ContextLogger) {
		ctx.LoggerName = LoggerName(name)
	}
}

// WithContextLevel specify logger level
func WithContextLevel(level int8) ContextLoggerOption {
	return func(ctx *ContextLogger) {
		ctx.Level = Level(level)
	}
}

// ContextLogger logger proxy used for specify logger name and level in std log adapter
type ContextLogger struct {
	LoggerName
	Level
}

func (ctx ContextLogger) Debug(msg string, keysAndValues ...interface{}) {
	GetLogger(ctx.LoggerName).Debug(msg, keysAndValues...)
}

func (ctx ContextLogger) Info(msg string, keysAndValues ...interface{}) {
	GetLogger(ctx.LoggerName).Info(msg, keysAndValues...)
}

func (ctx ContextLogger) Error(msg string, keysAndValues ...interface{}) {
	GetLogger(ctx.LoggerName).Error(msg, keysAndValues...)
}

func (ctx ContextLogger) String() string {
	return fmt.Sprintf("logger: %s: %s: ", ctx.LoggerName, ctx.Level)
}

var (
	_ LoggerInterface = (*ContextLogger)(nil)
)
