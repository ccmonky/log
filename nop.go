package log

// NopLogger is a logger which discards all logs.
type NopLogger struct{}

func (NopLogger) Error(msg string, keysAndValues ...interface{}) {}
func (NopLogger) Info(msg string, keysAndValues ...interface{})  {}
func (NopLogger) Debug(msg string, keysAndValues ...interface{}) {}

// Interface guard
var (
	_ LoggerInterface = (*NopLogger)(nil)
)
