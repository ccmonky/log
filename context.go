package log

import (
	"fmt"
)

type ContextLogger interface {
	LoggerInterface
	Level() Level
	Args() []interface{}
}

// Ctx returna a conetxt logger
// usage:
// - log.Ctx("xxx", log.ErrorLevel): use xxx logger in error level
// - log.Ctx("xxx"): use xxx logger in info level
// - log.Ctx(log.DebugLevel): use default logger in debug level
func Ctx(ps ...any) ContextLogger {
	var opts []contextLoggerOption
	for _, p := range ps {
		if p == nil {
			continue
		}
		switch p := p.(type) {
		case string:
			opts = append(opts, WithContextLoggerName(p))
		case LoggerName:
			opts = append(opts, WithContextLoggerName(string(p)))
		case int8:
			opts = append(opts, WithContextLevel(p))
		case Level:
			opts = append(opts, WithContextLevel(int8(p)))
		case contextLoggerOption:
			opts = append(opts, p)
		default:
			panic(fmt.Errorf("log.Ctx not support parameter type: %T", p))
		}
	}
	ctx := contextLogger{
		level: InfoLevel,
	}
	for _, opt := range opts {
		opt(&ctx)
	}
	return ctx
}

// contextLoggerOption specify context logger option
type contextLoggerOption func(*contextLogger)

// WithContextLoggerName specify logger name
func WithContextLoggerName(name string) contextLoggerOption {
	return func(ctx *contextLogger) {
		ctx.LoggerName = LoggerName(name)
	}
}

// WithContextLevel specify logger level
func WithContextLevel(level int8) contextLoggerOption {
	return func(ctx *contextLogger) {
		ctx.level = Level(level)
	}
}

// WithContextArgs specify logger level
func WithContextArgs(args []interface{}) contextLoggerOption {
	return func(ctx *contextLogger) {
		ctx.args = args
	}
}

// contextLogger logger proxy used for specify logger name and level in std log adapter
type contextLogger struct {
	LoggerName
	level Level
	args  []interface{}
}

func (ctx contextLogger) Debug(msg string, keysAndValues ...interface{}) {

	GetLogger(ctx.LoggerName).Debug(msg, append(ctx.args, keysAndValues...)...)
}

func (ctx contextLogger) Info(msg string, keysAndValues ...interface{}) {
	GetLogger(ctx.LoggerName).Info(msg, append(ctx.args, keysAndValues...)...)
}

func (ctx contextLogger) Error(msg string, keysAndValues ...interface{}) {
	GetLogger(ctx.LoggerName).Error(msg, append(ctx.args, keysAndValues...)...)
}

func (ctx contextLogger) String() string {
	return fmt.Sprintf("logger: %s: %s: ", ctx.LoggerName, ctx.level)
}

func (ctx contextLogger) Level() Level {
	return ctx.level
}

func (ctx contextLogger) Args() []interface{} {
	return nil
}

var (
	_ ContextLogger = (*contextLogger)(nil)
)
