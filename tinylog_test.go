package tinylog

import "testing"

func Test_tinylog(t *testing.T) {
	logger := NewStreamLogger(INFO)
	// debug msg will not show because level is INFO
	logger.Debug("this is [%s] msg.", "debug")
	logger.Info("this is [%s] msg.", "info")
	logger.Warn("this is [%s] msg!", "warn")
	logger.Error("this is [%s] msg!", "error")
	// fatal will print stacktrace and exit
	logger.Fatal("this is [%s] msg!", "fatal")
}
