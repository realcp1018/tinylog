## Introduction
*tinylog*: a simple logger based on golang built-in log

**Download:**
```
go get github.com/realcp1018/tinylog
```
**Usage**:

```
package main

import "github.com/realcp1018/tinylog"

func main() {
    StreamLogger := tinylog.NewStreamLogger(tinylog.INFO) 
    StreamLogger.Warn("warn msg")
    
    fileLogger := tinylog.NewFileLogger("test.log", tinylog.INFO)
    fileLogger.Warn("warn msg")
    
    // if you need some customized config: maxSizeMb, maxBackupCount, maxKeepDays
    fileLogger := tinylog.NewFileLogger("test.log", tinylog.INFO)
    fileLogger.SetFileConfig(128, 10, 7)
    fileLogger.Warn("warn msg")
}
```
Screen & test.log output :

`
2022/06/29 12:46:09.759870 [Warn] [main.go:7] warn msg
`

`Error()` and `Fatal()` will print stacktrace, `Fatal()` will do os.exit(1)

If you need a faster structured log for ETL & visualization, use https://github.com/uber-go/zap 