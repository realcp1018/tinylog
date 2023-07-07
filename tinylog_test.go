package tinylog

import (
	"fmt"
	"testing"
)

func Test_tinylog(t *testing.T) {
	var logger *TinyLogger

	// test StreamLogger
	logger = NewStreamLogger(INFO)
	logger.Info("this is a %s log to stdout", "INFO")
	fmt.Println(logger.GetLevelName())
	// test FileLogger
	logger = NewFileLogger("111.log", INFO)
	logger.SetFileConfig("", 1, 1, 1)
	logger.Info("this is a %s log to file", "INFO")

	// test default logger
	Info("default log")
	// convert the default logger to a FileLogger
	SetFileConfig("", 1, 1, 1)
	Info("default log to tiny.log file")
}
