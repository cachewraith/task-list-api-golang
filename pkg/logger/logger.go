package logger

import (
	"fmt"
	"sync"
	"time"
)

// LogLevel represents the severity of a log entry
type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     LogLevel  `json:"level"`
	Message   string    `json:"message"`
	Context   string    `json:"context"`
	Function  string    `json:"function"`
}

// Notifier interface for external alerting
type Notifier interface {
	SendAlert(entry LogEntry)
}

// Logger handles logging with in-memory storage for viewing
type Logger struct {
	entries  []LogEntry
	mu       sync.RWMutex
	maxSize  int
	notifier Notifier
}

var (
	instance *Logger
	once     sync.Once
)

// GetLogger returns the singleton logger instance
func GetLogger() *Logger {
	once.Do(func() {
		instance = &Logger{
			entries: make([]LogEntry, 0),
			maxSize: 1000, // Keep last 1000 entries
		}
	})
	return instance
}

// SetNotifier sets the external notifier for alerts
func (l *Logger) SetNotifier(n Notifier) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.notifier = n
}

// log adds a log entry
func (l *Logger) log(level LogLevel, context, function, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Context:   context,
		Function:  function,
	}

	l.entries = append(l.entries, entry)

	// Remove old entries if exceeding max size
	if len(l.entries) > l.maxSize {
		l.entries = l.entries[len(l.entries)-l.maxSize:]
	}

	// Send to external notifier if configured (for ERROR/WARN levels)
	if l.notifier != nil && (level == ERROR || level == WARN) {
		l.notifier.SendAlert(entry)
	}

	// Also print to console
	fmt.Printf("[%s] %s | %s.%s | %s\n",
		entry.Timestamp.Format("2006-01-02 15:04:05"),
		level,
		context,
		function,
		message,
	)
}

// Debug logs a debug message
func (l *Logger) Debug(context, function, message string) {
	l.log(DEBUG, context, function, message)
}

// Info logs an info message
func (l *Logger) Info(context, function, message string) {
	l.log(INFO, context, function, message)
}

// Warn logs a warning message
func (l *Logger) Warn(context, function, message string) {
	l.log(WARN, context, function, message)
}

// Error logs an error message
func (l *Logger) Error(context, function, message string) {
	l.log(ERROR, context, function, message)
}

// GetEntries returns all log entries (newest first)
func (l *Logger) GetEntries(limit int) []LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if limit <= 0 || limit > len(l.entries) {
		limit = len(l.entries)
	}

	// Return copy in reverse order (newest first)
	result := make([]LogEntry, limit)
	for i := 0; i < limit; i++ {
		result[i] = l.entries[len(l.entries)-1-i]
	}
	return result
}

// GetEntriesByLevel returns entries filtered by level
func (l *Logger) GetEntriesByLevel(level LogLevel, limit int) []LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var filtered []LogEntry
	for i := len(l.entries) - 1; i >= 0 && len(filtered) < limit; i-- {
		if l.entries[i].Level == level {
			filtered = append(filtered, l.entries[i])
		}
	}
	return filtered
}

// Clear clears all log entries
func (l *Logger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = make([]LogEntry, 0)
}

// Helper functions for easy access
func Debug(context, function, message string) {
	GetLogger().Debug(context, function, message)
}

func Info(context, function, message string) {
	GetLogger().Info(context, function, message)
}

func Warn(context, function, message string) {
	GetLogger().Warn(context, function, message)
}

func Error(context, function, message string) {
	GetLogger().Error(context, function, message)
}
