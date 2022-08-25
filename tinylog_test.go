package tinylog

import (
	"fmt"
	"testing"
)

func Test_tinylog(t *testing.T) {
	logger := NewStreamLogger(INFO)
	// debug msg will not show because level is INFO
	logger.Debug("this is [%s] msg.", "debug")
	logger.Info("this is [%s] msg.", "info")
	logger.Warn("this is [%s] msg!", "warn")
	logger.Error("this is [%s] msg!", "error")
	fmt.Printf("%v %T %v\n", logger.Flags(), logger.Writer(), logger.Prefix())
}
