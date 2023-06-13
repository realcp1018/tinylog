package tinylog

import (
	"fmt"
	"testing"
)

func Test_tinylog(t *testing.T) {
	var logger *TinyLogger

	logger = NewStreamLogger(INFO)
	logger.Info("this a %s log in stdout", "INFO")
	fmt.Println(logger.GetLevel())

	logger = NewFileLogger("111.log", INFO)
	logger.SetFileConfig(1, 1, 1)
	logger.Info("this a %s log in file", "INFO")
}
