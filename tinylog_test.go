package tinylog

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type failingWriter struct {
	err error
}

// Write returns the configured error for testing failed log writes.
func (w failingWriter) Write([]byte) (int, error) {
	return 0, w.err
}

// TestTinyLoggerLevelFiltering verifies that enabled levels are written and lower-priority levels are filtered.
func TestTinyLoggerLevelFiltering(t *testing.T) {
	var output bytes.Buffer
	logger := NewStreamLogger(ERROR)
	logger.setOutput(&output)

	logger.Print("filtered print")
	logger.Info("filtered info")
	logger.ErrorNoStackTrace("visible error")

	text := output.String()
	if strings.Contains(text, "filtered") {
		t.Fatalf("expected lower-priority logs to be filtered, got %q", text)
	}
	if !strings.Contains(text, "[ERROR]") || !strings.Contains(text, "visible error") {
		t.Fatalf("expected error log to be written, got %q", text)
	}
	if !strings.Contains(text, "tinylog_test.go:") {
		t.Fatalf("expected caller location in log, got %q", text)
	}
}

// TestTinyLoggerConcurrentLevelAccess verifies that level reads and writes can run with logging concurrently.
func TestTinyLoggerConcurrentLevelAccess(t *testing.T) {
	logger := NewStreamLogger(INFO)
	logger.setOutput(io.Discard)

	var waitGroup sync.WaitGroup
	for i := 0; i < 8; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			for j := 0; j < 100; j++ {
				logger.SetLevel(LogLevel(j % 5))
				_ = logger.GetLevelName()
				logger.Info("concurrent log")
			}
		}()
	}
	waitGroup.Wait()
}

// TestTinyLoggerWriteErrorHandler verifies that failed writes are reported to the configured handler.
func TestTinyLoggerWriteErrorHandler(t *testing.T) {
	expectedError := errors.New("write failed")
	logger := NewStreamLogger(INFO)
	logger.setOutput(failingWriter{err: expectedError})
	var reportedError error
	logger.SetWriteErrorHandler(func(err error) {
		reportedError = err
		logger.SetLevel(WARN)
	})

	done := make(chan struct{})
	go func() {
		logger.Info("failed log")
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("log write did not return")
	}
	if !errors.Is(reportedError, expectedError) {
		t.Fatalf("expected write error %v, got %v", expectedError, reportedError)
	}
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
	t.Cleanup(func() {
		_ = output.Close()
	})
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
