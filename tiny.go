// log.go is from:
// https://github.com/golang/go/blob/master/src/log/log.go
// I made some changes in the formatHeader() (to pretty the output), and delete(hide) some unreached functions/methods
// And defined some LogLevels && use lumberjack for log file management
package tinylog

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"runtime/debug"
	"strings"
	"sync"
)

type LogLevel uint

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type TinyLogger struct {
	*logger                 // go built-in logger
	mu           sync.Mutex // mutex add for SetPrefix
	logLevel     LogLevel
	filename     string // in case to reconfig lumberjack.Logger, we store filename here
	errorHandler func(error)
	callDepth    int
}

func NewFileLogger(fileName string, level LogLevel) *TinyLogger {
	logger := new(logger)
	logger.setOutput(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    512,
		MaxBackups: 16,
		MaxAge:     30,
		Compress:   true,
	})
	logger.setFlags(LstdFlags | Lmicroseconds | Lshortfile | Lmsgprefix)
	return &TinyLogger{
		logger:       logger,
		logLevel:     level,
		filename:     fileName,
		errorHandler: reportWriteError,
		callDepth:    2,
	}
}

func NewStreamLogger(level LogLevel) *TinyLogger {
	logger := new(logger)
	logger.setOutput(os.Stdout)
	logger.setFlags(LstdFlags | Lmicroseconds | Lmsgprefix | Lshortfile)
	return &TinyLogger{
		logger:       logger,
		logLevel:     level,
		errorHandler: reportWriteError,
		callDepth:    2,
	}
}

func newDefaultLogger(level LogLevel) *TinyLogger {
	logger := new(logger)
	logger.setOutput(os.Stdout)
	logger.setFlags(LstdFlags | Lmicroseconds | Lmsgprefix | Lshortfile)
	return &TinyLogger{
		logger:       logger,
		logLevel:     level,
		errorHandler: reportWriteError,
		callDepth:    3, // defaultLogger's log methods are wrapped by pkg functions, so callDepth+1
	}
}

// SetFileConfig set file configs for FileLogger, or convert a StreamLogger to FileLogger
func (l *TinyLogger) SetFileConfig(fileName string, maxSizeMb, maxBackupCount, maxKeepDays int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	oldWriter := l.writer()

	var newFileName string
	if fileName == "" {
		newFileName = l.filename
	} else {
		newFileName = fileName
	}
	// if fileName not given for StreamLogger, then set fileName to "tiny.log"
	if newFileName == "" {
		newFileName = "tiny.log"
	}
	newWriter := &lumberjack.Logger{
		Filename:   newFileName,
		MaxSize:    maxSizeMb,
		MaxBackups: maxBackupCount,
		MaxAge:     maxKeepDays,
		Compress:   true,
	}
	l.setOutput(newWriter)
	l.filename = newFileName
	if oldFile, ok := oldWriter.(*lumberjack.Logger); ok {
		_ = oldFile.Close()
	}
}

// SetWriteErrorHandler sets the callback for log write errors. A nil handler restores the default stderr reporter.
func (l *TinyLogger) SetWriteErrorHandler(handler func(error)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if handler == nil {
		l.errorHandler = reportWriteError
		return
	}
	l.errorHandler = handler
}

// write writes a formatted message when the requested level is enabled.
func (l *TinyLogger) write(level LogLevel, prefix, format string, v ...interface{}) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel > level {
		return false, nil
	}
	l.setPrefix(prefix)
	return true, l.output(l.callDepth+1, fmt.Sprintf(format, v...))
}

// handleWriteError reports a failed log write without holding the logger lock.
func (l *TinyLogger) handleWriteError(err error) {
	if err == nil {
		return
	}
	l.mu.Lock()
	handler := l.errorHandler
	l.mu.Unlock()
	if handler == nil {
		reportWriteError(err)
		return
	}
	handler(err)
}

func (l *TinyLogger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logLevel = level
}

func (l *TinyLogger) GetLevelName() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	switch l.logLevel {
	case 0:
		return "DEBUG"
	case 1:
		return "INFO"
	case 2:
		return "WARN"
	case 3:
		return "ERROR"
	case 4:
		return "FATAL"
	default:
		return ""
	}
}

func (l *TinyLogger) Debug(format string, v ...interface{}) {
	_, err := l.write(DEBUG, "[DEBUG] ", format, v...)
	l.handleWriteError(err)
}

func (l *TinyLogger) Info(format string, v ...interface{}) {
	_, err := l.write(INFO, "[INFO] ", format, v...)
	l.handleWriteError(err)
}

func (l *TinyLogger) Warn(format string, v ...interface{}) {
	_, err := l.write(WARN, "[WARN] ", format, v...)
	l.handleWriteError(err)
}

func (l *TinyLogger) Error(format string, v ...interface{}) {
	_, err := l.write(ERROR, "[ERROR] ", "%s\n[stacktrace]:\n%s", fmt.Sprintf(format, v...), string(debug.Stack()))
	l.handleWriteError(err)
}

// ErrorNoStackTrace print error with no stacktrace
func (l *TinyLogger) ErrorNoStackTrace(format string, v ...interface{}) {
	_, err := l.write(ERROR, "[ERROR] ", format, v...)
	l.handleWriteError(err)
}

// Fatal do exit
func (l *TinyLogger) Fatal(format string, v ...interface{}) {
	written, err := l.write(FATAL, "[FATAL] ", "%s\n[stacktrace]:\n%s", fmt.Sprintf(format, v...), string(debug.Stack()))
	l.handleWriteError(err)
	if written {
		os.Exit(1)
	}
}

// Add some common print functions(for interfaces)
func (l *TinyLogger) Print(v ...interface{}) {
	var format string
	for i := 0; i < len(v); i++ {
		format += "%v "
	}
	_, err := l.write(WARN, "[WARN] ", strings.TrimSpace(format), v...)
	l.handleWriteError(err)
}

func (l *TinyLogger) Printf(format string, v ...interface{}) {
	_, err := l.write(WARN, "[WARN] ", format, v...)
	l.handleWriteError(err)
}

func (l *TinyLogger) Println(v ...interface{}) {
	var format string
	for i := 0; i < len(v); i++ {
		format += "%v "
	}
	_, err := l.write(WARN, "[WARN] ", strings.TrimSpace(format)+"\n", v...)
	l.handleWriteError(err)
}

// reportWriteError reports log write failures to stderr.
func reportWriteError(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "tinylog: failed to write log: %v\n", err)
}
