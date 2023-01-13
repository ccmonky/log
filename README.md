# log

log interface and default implementation for ccmonky.

## Usage

### inject logger

you can inject logger with the following interface(e.g. zap.SugaredLogger):

```go
type LoggerInterface interface {
    Debug(msg string, keysAndValues ...interface{})
    Info(msg string, keysAndValues ...interface{})
    Error(msg string, keysAndValues ...interface{})
}
```

inject logger like this:

```go
log.SetLogger("xxx", zap.S())
```

### log levels

ccmonky/log support these log levels by default:

- ErrorLevel: error level, also can use short alias `E`
- InfoLevel: error level, also can use short alias `I`
- DebugLevel: error level, also can use short alias `D`

### std log adapter

ccmonky/log has std log adapter support, you can use it like this:

- first, replace your std log import to `github.com/ccmonky/log`

that's all for existing libs, and then in app, replace the default logger like this:

```go
// NOTE: 
// 1. default logger name is ""
// 2. appLogger is your app logger, which should implement `log.LoggerInterface`
log.SetLogger("", appLogger) 
```

then all your logs of libs will be recorded by appLogger, but note that, all the logs can only be recorded by
appLogger in info level.

If you are develop new libs, you can also support logger and levels like this:

```go
log.Println(log.Ctx("xxx", log.E), ...)
log.Printf("%v %d %s", log.Ctx("yyy", log.D), code, msg)
```
