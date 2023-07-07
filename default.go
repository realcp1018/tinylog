package tinylog

// default logger
var defaultLogger = NewStreamLogger(INFO)

func SetFileConfig(fileName string, maxSizeMb, maxBackupCount, maxKeepDays int) {
	defaultLogger.SetFileConfig(fileName, maxSizeMb, maxBackupCount, maxKeepDays)
}

func SetLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}

func GetLevelName() string {
	return defaultLogger.GetLevelName()
}

func Debug(format string, v ...interface{}) {
	defaultLogger.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	defaultLogger.Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	defaultLogger.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	defaultLogger.Error(format, v...)
}

func ErrorNoStackTrace(format string, v ...interface{}) {
	defaultLogger.ErrorNoStackTrace(format, v...)
}

func Fatal(format string, v ...interface{}) {
	defaultLogger.Fatal(format, v...)
}
