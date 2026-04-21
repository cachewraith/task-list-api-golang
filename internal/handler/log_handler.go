package handler

import (
	"net/http"
	"strconv"

	"github.com/example/todo-api/pkg/logger"
	"github.com/example/todo-api/pkg/response"
	"github.com/go-chi/chi/v5"
)

// LogHandler handles log viewing requests
type LogHandler struct{}

// NewLogHandler creates a new log handler
func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

// Routes returns the router with all log routes registered
func (h *LogHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetLogs)
	r.Delete("/", h.ClearLogs)
	return r
}

// GetLogs returns recent log entries
func (h *LogHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	// Get limit from query param (default 100)
	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	// Get level filter from query param
	level := r.URL.Query().Get("level")

	var entries []logger.LogEntry
	if level != "" {
		entries = logger.GetLogger().GetEntriesByLevel(logger.LogLevel(level), limit)
	} else {
		entries = logger.GetLogger().GetEntries(limit)
	}

	logger.Info("LOG_HANDLER", "GetLogs", "Retrieved "+strconv.Itoa(len(entries))+" log entries")

	response.OK(w, map[string]interface{}{
		"entries": entries,
		"total":   len(entries),
		"limit":   limit,
	})
}

// ClearLogs clears all log entries
func (h *LogHandler) ClearLogs(w http.ResponseWriter, r *http.Request) {
	logger.GetLogger().Clear()
	logger.Info("LOG_HANDLER", "ClearLogs", "All logs cleared")

	response.NewSuccessWithMessage("All logs cleared", nil).Write(w, http.StatusOK)
}
