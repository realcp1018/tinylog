package tinylog

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
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

// TestTinyLoggerSetFileConfigClosesPreviousFile verifies that reconfiguration closes the old log file.
func TestTinyLoggerSetFileConfigClosesPreviousFile(t *testing.T) {
	logDir := t.TempDir()
	oldFileName := filepath.Join(logDir, "old.log")
	newFileName := filepath.Join(logDir, "new.log")
	logger := NewFileLogger(oldFileName, INFO)
	oldWriter, ok := logger.writer().(*lumberjack.Logger)
	if !ok {
		t.Fatalf("expected lumberjack output, got %T", logger.writer())
	}
	t.Cleanup(func() {
		_ = oldWriter.Close()
		if newWriter, ok := logger.writer().(*lumberjack.Logger); ok {
			_ = newWriter.Close()
		}
	})

	logger.Info("write to old file")
	logger.SetFileConfig(newFileName, 1, 1, 1)
	if err := os.Remove(oldFileName); err != nil {
		t.Fatalf("expected old log file to be closed: %v", err)
	}
}
