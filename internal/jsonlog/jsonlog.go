package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// Level represents the severity level for log entry.
type Level int8

// The constants represent a specific severity level.
const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
)

// String method return a human-friendly string for severity level.
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// Logger is a custom logger that holds the output destination
// that the log entries will be written to, the minimum severity
// level that logs entries will be written for, and mutex for coordinating writes.
type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

// New return a new Logger instance which writes log entries at or
// above a minimum severity level to a specific destination.
func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

// PrintInfo logs the message and properties if present to a specific destination
// with LevelInfo assigned as the severity level.
func (l *Logger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

// PrintError logs the error and properties if present to a specific destination
// with LevelError assigned as the severity level.
func (l *Logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

// PrintFatal logs the error and properties if present to a specific destination
// with LevelInfo assigned as the severity level and also terminate the app
// with exit code 1.
func (l *Logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

// print is an internal method for writing the log entry.
func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	// include the stack trace for entries at the Error and Fatal levels
	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte
	line, err := json.Marshal(aux)

	// Marshall the anonymous struct to JSON and store it in the line variable.
	// If there was a problem creating the JSON, set the contents of the log
	// entry to be the plain-text error message instead.
	if err != nil {
		line = []byte(err.Error() + ": unable to marshal log message")
	}

	// Lock the mutex so that no two writes to the output destination can
	// happen concurrently.
	l.mu.Lock()
	defer l.mu.Unlock()

	// write the log entry followed by a new line
	return l.out.Write(append(line, '\n'))
}

// Write implements the io.Writer interface, by writing the log entry
// with LevelError as the severity level.
func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
