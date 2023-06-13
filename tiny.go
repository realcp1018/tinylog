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
	"strings"
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
	// mutex add for SetPrefix
	mu       sync.Mutex
	logLevel uint8
	// in case to reconfig lumberjack.Logger, we store filename here
	filename string
}

func NewFileLogger(fileName string, level uint8) *TinyLogger {
	logger := new(Logger)
	logger.SetOutput(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    512,
		MaxBackups: 16,
		MaxAge:     30,
		Compress:   true,
	})
	logger.SetFlags(LstdFlags | Lmicroseconds | Lshortfile | Lmsgprefix)
	return &TinyLogger{
		Logger:   logger,
		logLevel: level,
		filename: fileName,
	}
}

func NewStreamLogger(level uint8) *TinyLogger {
	logger := new(Logger)
	logger.SetOutput(os.Stdout)
	logger.SetFlags(LstdFlags | Lmicroseconds | Lmsgprefix | Lshortfile)
	return &TinyLogger{
		Logger:   logger,
		logLevel: level,
	}
}

func (l *TinyLogger) SetFileConfig(maxSizeMb, maxBackupCount, maxKeepDays int) {
	switch l.Writer().(type) {
	case *lumberjack.Logger:
		l.SetOutput(&lumberjack.Logger{
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

func (l *TinyLogger) SetLevel(level uint8) {
	l.logLevel = level
}

func (l *TinyLogger) GetLevel() uint8 {
	return l.logLevel
}

func (l *TinyLogger) Debug(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel == DEBUG {
		l.SetPrefix(fmt.Sprintf("[DEBUG] "))
		_ = l.Output(2, fmt.Sprintf(format, v...))
	}
}

func (l *TinyLogger) Info(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel <= INFO {
		l.SetPrefix("[INFO] ")
		_ = l.Output(2, fmt.Sprintf(format, v...))
	}
}

func (l *TinyLogger) Warn(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel <= WARN {
		l.SetPrefix("[WARN] ")
		_ = l.Output(2, fmt.Sprintf(format, v...))
	}
}

func (l *TinyLogger) Error(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel <= ERROR {
		l.SetPrefix("[ERROR] ")
		_ = l.Output(2, fmt.Sprintf("%s\n[stacktrace]:\n%s", fmt.Sprintf(format, v...), string(debug.Stack())))
	}
}

// Fatal do exit
func (l *TinyLogger) Fatal(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel <= FATAL {
		l.SetPrefix("[FATAL] ")
		_ = l.Output(2, fmt.Sprintf("%s\n[stacktrace]:\n%s", fmt.Sprintf(format, v...), string(debug.Stack())))
		os.Exit(1)
	}
}

// Add some common print interface functions
func (l *TinyLogger) Print(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logLevel <= INFO {
		l.SetPrefix("[INFO] ")
		var format string
		for i := 0; i < len(v); i++ {
			format += "%v "
		}
		_ = l.Output(2, fmt.Sprintf(strings.TrimSpace(format), v...))
	}
}

func (l *TinyLogger) Printf(format string, v ...interface{}) {
	l.Info(format, v...)
}

func (l *TinyLogger) Println(v ...interface{}) {
	l.Print(v, "\n")
}
