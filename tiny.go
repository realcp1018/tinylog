// log.go is a copy of stdlib log from:
// https://github.com/golang/go/blob/master/src/log/log.go
// but I made some changes to stdlib log in the formatHeader() function(to pretty the output)
// this file(tiny.go) add some log levels to stdlib log and hardcoded some configs of lumberjack & logger flags
package tinylog

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"runtime/debug"
	"sync"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

type TinyLogger struct {
	*Logger
	mu       sync.Mutex
	logLevel uint8
}

func NewFileLogger(fileName string, level uint8) *TinyLogger {
	logger := new(Logger)
	logger.SetOutput(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    500,
		MaxBackups: 30,
		MaxAge:     7,
		Compress:   true,
	})
	logger.SetFlags(LstdFlags | Lmicroseconds | Lshortfile | Lmsgprefix)
	return &TinyLogger{Logger: logger, logLevel: level}
}

func NewStreamLogger(level uint8) *TinyLogger {
	logger := new(Logger)
	logger.SetOutput(os.Stdout)
	logger.SetFlags(LstdFlags | Lmicroseconds | Lmsgprefix | Lshortfile)
	return &TinyLogger{Logger: logger, logLevel: level}
}

func (l *TinyLogger) SetLevel(level uint8) {
	l.logLevel = level
}

func (l *TinyLogger) GetLevel() uint8 {
	return l.logLevel
}

func (l *TinyLogger) Debug(logStr string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel == DEBUG {
		l.SetPrefix(fmt.Sprintf("[DEBUG] "))
		_ = l.Output(2, fmt.Sprintf(logStr, v...))
	}
}

func (l *TinyLogger) Info(logStr string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel <= INFO {
		l.SetPrefix("[INFO] ")
		_ = l.Output(2, fmt.Sprintf(logStr, v...))
	}
}

func (l *TinyLogger) Warn(logStr string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel <= WARN {
		l.SetPrefix( "[WARN] ")
		_ = l.Output(2, fmt.Sprintf(logStr, v...))
	}
}

func (l *TinyLogger) Error(logStr string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel <= ERROR {
		l.SetPrefix("[ERROR] ")
		_ = l.Output(2, fmt.Sprintf(logStr, v...))
	}
}

// Fatal do exit
func (l *TinyLogger) Fatal(logStr string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel <= FATAL {
		l.SetPrefix("[FATAL] ")
		_ = l.Output(2, fmt.Sprintf("%s [stacktrace]:\n%s", fmt.Sprintf(logStr, v...), string(debug.Stack())))
		os.Exit(1)
	}
}
