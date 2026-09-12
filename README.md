## Introduction
*tinylog*: a simple logger based on golang built-in log, thread safe.

#### Download:
```
go get github.com/realcp1018/tinylog
```
#### Usage:

```
package main

import "github.com/realcp1018/tinylog"

func main() {
    streamLogger := tinylog.NewStreamLogger(tinylog.INFO)
    streamLogger.Warn("warn msg")

    fileLogger := tinylog.NewFileLogger("test.log", tinylog.INFO)
    fileLogger.Warn("warn msg")

    // if you need some customized config: maxSizeMb, maxBackupCount, maxKeepDays
    fileLogger.SetFileConfig("", 128, 10, 7)
    fileLogger.Warn("warn msg")
}
```
Screen & test.log output :

```
2022/06/29 12:46:09.759870 [WARN] [main.go:7] warn msg
```

`Error()` and `Fatal()` will print stacktrace, `Fatal()` will do os.exit(1).

#### HowTo:
1. how to use it in a project?
```go
// create log package with file log.go
package logger

import "github.com/realcp1018/tinylog"

var GlobalLogger = tinylog.NewFileLogger("<log file path>", tinylog.INFO)

// then use it in other packages
//...
logger.GlobalLogger.Info("this is info msg")
//...
```
2. can I use tinylog without create a GlobalLogger?

Yes, you can use the default logger in tinylog:
```go
import "github.com/realcp1018/tinylog"

...
// default logger default write to os.stdout 
tinylog.Info("default log")

// convert the default logger to a FileLogger
tinylog.SetFileConfig("mylog.log", 1, 1, 1)
tinylog.Info("default log to mylog.log file")
...
```
Once you converted the default logger to a FileLogger, logs will write to logfile anywhere else in your project,
because tinylog will only compile once.

By default, tinylog reports log write errors to stderr. You can provide a custom handler:
```go
tinylog.SetWriteErrorHandler(func(err error) {
    // send err to your monitoring system
})
```
#### More:

If you need a faster structured log for ETL & visualization, use https://github.com/uber-go/zap
