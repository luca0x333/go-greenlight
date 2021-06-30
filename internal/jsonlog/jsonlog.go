package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// Level custom type define the severity level for a log entry.
type Level int8

const (
	LevelInfo  Level = iota // Has the value 0
	LevelError              // Has the value 1
	LevelFatal              // Has the value 2
	LevelOff                // Has the value 3
)

// String returns a human friendly string for the severity level.
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

// Logger struct holds the output destination that the log entries will be written to,
// the minimum log severity level that log entries will be written for and
// a mutex for coordinating writes.
type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

// New function return a new Logger instance.
func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

// Print is an internal method for writing the log entry.
func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	// Declare an anonymous struct holding the data for the log entry.
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

	// Include a stack trace for entries at ERROR and FATAL levels.
	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	// Marshal the anonymous struct to JSON and store it in the line variable.
	// If there is any error set the contents of the log entry as plain-text error message.
	var line []byte
	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message:" + err.Error())
	}

	// Lock the mutex so that no two writes ad the destination can happen concurrently.
	// Without this mutex lock, it’s possible that the content of multiple log entries
	// would be written at exactly the same time and be mixed up in the output,
	// rather than each entry being written in full on its own line.
	l.mu.Lock()
	defer l.mu.Unlock()

	// Write the log entry followed by a new line.
	return l.out.Write(append(line, '\n'))
}

func (l *Logger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *Logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

// We also implement a Write() method on our Logger type so that it satisfies the io.Writer interface
// This writes a log entry at the ERROR level with no additional properties.
func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
