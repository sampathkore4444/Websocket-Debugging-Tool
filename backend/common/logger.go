package common

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Logger struct {
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
	debug *log.Logger
}

var logger *Logger

func InitLogger(format string, output string) *Logger {
	var writer io.Writer

	if output == "file" {
		f, err := os.OpenFile("wsinspect.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		writer = f
	} else {
		writer = os.Stdout
	}

	logger = &Logger{
		info:  log.New(writer, "[INFO] ", log.LstdFlags),
		warn:  log.New(writer, "[WARN] ", log.LstdFlags),
		error: log.New(writer, "[ERROR] ", log.LstdFlags),
		debug: log.New(writer, "[DEBUG] ", log.LstdFlags),
	}

	return logger
}

func GetLogger() *Logger {
	if logger == nil {
		return InitLogger("text", "stdout")
	}
	return logger
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.info.Output(2, fmt.Sprintf(format, args...))
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.warn.Output(2, fmt.Sprintf(format, args...))
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.error.Output(2, fmt.Sprintf(format, args...))
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.debug.Output(2, fmt.Sprintf(format, args...))
}

// RequestLogger middleware
func RequestLogger(next func(format string, args ...interface{})) func(format string, args ...interface{}) {
	return func(format string, args ...interface{}) {
		start := time.Now()
		next(format, args...)
		duration := time.Since(start)
		GetLogger().Info("%s %v", fmt.Sprintf(format, args...), duration)
	}
}
