package transcoder

import (
	"log"
	"os"
)

// LogLevel уровень логирования
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// Logger интерфейс для логирования
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	SetLevel(level LogLevel)
}

// DefaultLogger стандартная реализация логгера
type DefaultLogger struct {
	level  LogLevel
	logger *log.Logger
}

// NewDefaultLogger создает новый логгер
func NewDefaultLogger(level LogLevel) *DefaultLogger {
	return &DefaultLogger{
		level:  level,
		logger: log.New(os.Stdout, "[TRANSCODER] ", log.LstdFlags),
	}
}

func (l *DefaultLogger) Debug(msg string, args ...interface{}) {
	if l.level <= LogLevelDebug {
		l.logger.Printf("[DEBUG] "+msg, args...)
	}
}

func (l *DefaultLogger) Info(msg string, args ...interface{}) {
	if l.level <= LogLevelInfo {
		l.logger.Printf("[INFO] "+msg, args...)
	}
}

func (l *DefaultLogger) Warn(msg string, args ...interface{}) {
	if l.level <= LogLevelWarn {
		l.logger.Printf("[WARN] "+msg, args...)
	}
}

func (l *DefaultLogger) Error(msg string, args ...interface{}) {
	if l.level <= LogLevelError {
		l.logger.Printf("[ERROR] "+msg, args...)
	}
}

func (l *DefaultLogger) SetLevel(level LogLevel) {
	l.level = level
}

// NoOpLogger логгер который ничего не делает
type NoOpLogger struct{}

func (NoOpLogger) Debug(msg string, args ...interface{}) {}
func (NoOpLogger) Info(msg string, args ...interface{})  {}
func (NoOpLogger) Warn(msg string, args ...interface{})  {}
func (NoOpLogger) Error(msg string, args ...interface{}) {}
func (NoOpLogger) SetLevel(level LogLevel)               {}
