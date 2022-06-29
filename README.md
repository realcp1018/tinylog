A tiny log based on go stdlib log:
```
logger := NewStreamLogger(INFO)
logger.Info("this is [%s] msg.", "info")
```
Output:

`
2022/06/29 11:59:51.842839 [Info] [tinylog_test.go:8] This is [info] msg.
`