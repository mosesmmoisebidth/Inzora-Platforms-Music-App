package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger interface for structured logging
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	With(fields ...interface{}) Logger
}

type logger struct {
	entry *logrus.Entry
}

// NewLogger creates a new structured logger
func NewLogger() Logger {
	log := logrus.New()
	
	// Set format based on environment
	env := strings.ToLower(os.Getenv("APP_ENVIRONMENT"))
	if env == "production" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		})
	}

	// Set log level
	level := strings.ToLower(os.Getenv("APP_LOG_LEVEL"))
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	return &logger{
		entry: logrus.NewEntry(log),
	}
}

// NewLoggerWithFields creates a logger with predefined fields
func NewLoggerWithFields(fields map[string]interface{}) Logger {
	log := NewLogger().(*logger)
	return &logger{
		entry: log.entry.WithFields(logrus.Fields(fields)),
	}
}

func (l *logger) Debug(msg string, fields ...interface{}) {
	l.entry.WithFields(parseFields(fields...)).Debug(msg)
}

func (l *logger) Info(msg string, fields ...interface{}) {
	l.entry.WithFields(parseFields(fields...)).Info(msg)
}

func (l *logger) Warn(msg string, fields ...interface{}) {
	l.entry.WithFields(parseFields(fields...)).Warn(msg)
}

func (l *logger) Error(msg string, fields ...interface{}) {
	l.entry.WithFields(parseFields(fields...)).Error(msg)
}

func (l *logger) Fatal(msg string, fields ...interface{}) {
	l.entry.WithFields(parseFields(fields...)).Fatal(msg)
}

func (l *logger) With(fields ...interface{}) Logger {
	return &logger{
		entry: l.entry.WithFields(parseFields(fields...)),
	}
}

// parseFields converts key-value pairs to logrus.Fields
func parseFields(fields ...interface{}) logrus.Fields {
	if len(fields) == 0 {
		return logrus.Fields{}
	}

	// If the first field is already a map, use it
	if f, ok := fields[0].(map[string]interface{}); ok {
		return logrus.Fields(f)
	}

	// Parse key-value pairs
	result := logrus.Fields{}
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			result[key] = fields[i+1]
		}
	}

	return result
}
