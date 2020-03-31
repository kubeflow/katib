package goptuna

import (
	"fmt"
	"log"
)

// Logger is the interface for logging messages.
// If you want to log nothing, please set Logger as nil.
// If you want to print more verbose logs,
// it might StudyOptionSetTrialNotifyChannel option are useful.
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

// LoggerLevel represents Level is a logging priority.
// Higher levels are more important.
type LoggerLevel int

const (
	// LoggerLevelDebug logs are typically voluminous, and are usually disabled in production.
	LoggerLevelDebug LoggerLevel = iota - 1
	// LoggerLevelInfo is the default logging priority.
	LoggerLevelInfo
	// LoggerLevelWarn logs are more important than Info, but don't need individual human review.
	LoggerLevelWarn
	// LoggerLevelError logs are high-priority.
	LoggerLevelError
)

var _ Logger = &StdLogger{}

// StdLogger wraps 'log' standard library.
type StdLogger struct {
	Logger *log.Logger
	Level  LoggerLevel
	Color  bool
}

func (l *StdLogger) write(msg string, fields ...interface{}) {
	if l.Logger == nil {
		return
	}
	fields = append([]interface{}{msg}, fields...)
	l.Logger.Println(fields...)
}

// Debug logs a message at LoggerLevelDebug which are usually disabled in production.
func (l *StdLogger) Debug(msg string, fields ...interface{}) {
	if l.Level > LoggerLevelDebug {
		return
	}

	prefix := "[DEBUG] "
	if l.Color {
		prefix = "\033[1;36m" + prefix + "\033[0m"
	}
	l.write(fmt.Sprintf("%s%s", prefix, msg), fields...)
}

// Info logs a message at LoggerLevelInfo. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
func (l *StdLogger) Info(msg string, fields ...interface{}) {
	if l.Level > LoggerLevelInfo {
		return
	}

	prefix := "[INFO] "
	if l.Color {
		prefix = "\033[1;34m" + prefix + "\033[0m"
	}
	l.write(fmt.Sprintf("%s%s:", prefix, msg), fields...)
}

// Warn logs a message at LoggerLevelWarn which more important than Info,
// but don't need individual human review.
func (l *StdLogger) Warn(msg string, fields ...interface{}) {
	if l.Level > LoggerLevelWarn {
		return
	}

	prefix := "[WARN] "
	if l.Color {
		prefix = "\033[1;33m" + prefix + "\033[0m"
	}
	l.write(fmt.Sprintf("%s%s:", prefix, msg), fields...)
}

// Error logs a message at LoggerLevelError which are high-priority.
func (l *StdLogger) Error(msg string, fields ...interface{}) {
	if l.Level > LoggerLevelError {
		return
	}

	prefix := "[ERROR] "
	if l.Color {
		prefix = "\033[1;31m" + prefix + "\033[0m"
	}
	l.write(fmt.Sprintf("%s%s:", prefix, msg), fields...)
}
