package tinylog

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
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
	Info("default log[default logger]")
	// convert the default logger to a FileLogger
	SetFileConfig("", 1, 1, 1)
	Info("default log to tiny.log file[default logger]")
}

// TestTinyLoggerSetFileConfigUpdatesFilename verifies that an empty filename keeps the latest configured path.
func TestTinyLoggerSetFileConfigUpdatesFilename(t *testing.T) {
	currentFile := "current.log"

	logger := NewFileLogger("original.log", INFO)
	logger.SetFileConfig(currentFile, 1, 1, 1)
	logger.SetFileConfig("", 1, 1, 1)

	output, ok := logger.writer().(*lumberjack.Logger)
	if !ok {
		t.Fatalf("expected lumberjack output, got %T", logger.writer())
	}
	if output.Filename != currentFile {
		t.Fatalf("expected current log file %q, got %q", currentFile, output.Filename)
	}
}
