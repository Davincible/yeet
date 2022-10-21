package logger

import (
	"http-poc/logger/zerolog"
)

var _ Logger = (*zerolog.ZeroLogger)(nil)

type Level int8

type Logger interface {
	// Fields set fields to always be logged
	// Fields(fields map[string]interface{}) Logger
	// Log writes a log entry
	// Log(level Level, v ...interface{})
	// // Logf writes a formatted log entry
	// Logf(level Level, format string, v ...interface{})

	Info(...any)
	Infof(string, ...any)

	Warning(...any)
	Warningf(string, ...any)

	Error(...any)
	Errorf(string, ...any)

	Fatal(...any)
	Fatalf(string, ...any)

	Debug(...any)
	Debugf(string, ...any)

	Trace(...any)
	Tracef(string, ...any)

	// String returns the name of logger
	String() string
}
