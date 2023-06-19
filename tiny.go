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
	*logger             // go built-in logger
	mu       sync.Mutex // mutex add for SetPrefix
	logLevel LogLevel
	filename string // in case to reconfig lumberjack.Logger, we store filename here
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
		logger:   logger,
		logLevel: level,
		filename: fileName,
	}
}

func NewStreamLogger(level LogLevel) *TinyLogger {
	logger := new(logger)
	logger.setOutput(os.Stdout)
	logger.setFlags(LstdFlags | Lmicroseconds | Lmsgprefix | Lshortfile)
	return &TinyLogger{
		logger:   logger,
		logLevel: level,
	}
}

func (l *TinyLogger) SetFileConfig(maxSizeMb, maxBackupCount, maxKeepDays int) {
	switch l.writer().(type) {
	case *lumberjack.Logger:
		l.setOutput(&lumberjack.Logger{
			Filename:   l.filename,
			MaxSize:    maxSizeMb,
			MaxBackups: maxBackupCount,
			MaxAge:     maxKeepDays,
			Compress:   true,
		})
	default:
		return
	}
}

func (l *TinyLogger) SetLevel(level LogLevel) {
	l.logLevel = level
}

func (l *TinyLogger) GetLevelName() string {
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
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel == DEBUG {
		l.setPrefix(fmt.Sprintf("[DEBUG] "))
		_ = l.output(2, fmt.Sprintf(format, v...))
	}
}

func (l *TinyLogger) Info(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel <= INFO {
		l.setPrefix("[INFO] ")
		_ = l.output(2, fmt.Sprintf(format, v...))
	}
}

func (l *TinyLogger) Warn(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel <= WARN {
		l.setPrefix("[WARN] ")
		_ = l.output(2, fmt.Sprintf(format, v...))
	}
}

func (l *TinyLogger) Error(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel <= ERROR {
		l.setPrefix("[ERROR] ")
		_ = l.output(2, fmt.Sprintf("%s\n[stacktrace]:\n%s", fmt.Sprintf(format, v...), string(debug.Stack())))
	}
}

// ErrorNoStackTrace print error with no stacktrace
func (l *TinyLogger) ErrorNoStackTrace(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel <= ERROR {
		l.setPrefix("[ERROR] ")
		_ = l.output(2, fmt.Sprintf(format, v...))
	}
}

// Fatal do exit
func (l *TinyLogger) Fatal(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel <= FATAL {
		l.setPrefix("[FATAL] ")
		_ = l.output(2, fmt.Sprintf("%s\n[stacktrace]:\n%s", fmt.Sprintf(format, v...), string(debug.Stack())))
		os.Exit(1)
	}
}

// Add some common print functions(for interfaces)
func (l *TinyLogger) Print(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel <= INFO {
		l.setPrefix("[INFO] ")
		var format string
		for i := 0; i < len(v); i++ {
			format += "%v "
		}
		_ = l.output(2, fmt.Sprintf(strings.TrimSpace(format), v...))
	}
}

func (l *TinyLogger) Printf(format string, v ...interface{}) {
	l.Info(format, v...)
}

func (l *TinyLogger) Println(v ...interface{}) {
	l.Print(v, "\n")
}
