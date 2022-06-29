A tiny log based on go stdlib log.

```
go get github.com/realcp1018/tinylog
```
Usage:

```
package main

import "github.com/realcp1018/tinylog"

func main() {
    logger := tinylog.NewStreamLogger(tinylog.INFO) 
    logger.Warn("warn msg")
}
```
Output:

`
2022/06/29 12:46:09.759870 [Warn] [main.go:7] warn msg
`