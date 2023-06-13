package tinylog

import (
	"testing"
)

func Test_tinylog(t *testing.T) {
	var logger *TinyLogger

	logger = NewStreamLogger(INFO)
	logger.Info("this a %s log in stdout", "INFO")

	logger = NewFileLogger("111.log", INFO)
	logger.SetFileConfig(1, 1, 1)
	logger.Info("this a %s log in file", "INFO")
}
