package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cachewraith/task-list-api-golang/pkg/logger"
)

// TelegramNotifier sends log alerts to Telegram
type TelegramNotifier struct {
	botToken  string
	groupID   string
	threadID  string
	client    *http.Client
	enabled   bool
	mu        sync.RWMutex
	minLevel  logger.LogLevel
}

var (
	instance *TelegramNotifier
	once     sync.Once
)

// GetTelegramNotifier returns singleton instance
func GetTelegramNotifier() *TelegramNotifier {
	once.Do(func() {
		instance = &TelegramNotifier{
			client: &http.Client{
				Timeout: 5 * time.Second,
			},
			minLevel: logger.ERROR, // Default: only ERROR and above
		}
	})
	return instance
}

// Configure sets up the notifier with credentials
func (t *TelegramNotifier) Configure(botToken, groupID, threadID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	t.botToken = botToken
	t.groupID = groupID
	t.threadID = threadID
	t.enabled = botToken != "" && groupID != ""
}

// SetMinLevel sets minimum log level to send alerts
func (t *TelegramNotifier) SetMinLevel(level logger.LogLevel) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.minLevel = level
}

// shouldNotify checks if log level meets minimum threshold
func (t *TelegramNotifier) shouldNotify(level logger.LogLevel) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	
	if !t.enabled {
		return false
	}
	
	levels := map[logger.LogLevel]int{
		logger.DEBUG: 0,
		logger.INFO:  1,
		logger.WARN:  2,
		logger.ERROR: 3,
	}
	
	return levels[level] >= levels[t.minLevel]
}

// SendAlert sends log entry to Telegram (non-blocking)
func (t *TelegramNotifier) SendAlert(entry logger.LogEntry) {
	if !t.shouldNotify(entry.Level) {
		return
	}
	
	// Run async to not block the main application
	go t.send(entry)
}

// send performs the actual HTTP request to Telegram API
func (t *TelegramNotifier) send(entry logger.LogEntry) {
	t.mu.RLock()
	botToken := t.botToken
	groupID := t.groupID
	threadID := t.threadID
	t.mu.RUnlock()
	
	message := t.formatMessage(entry)
	
	payload := map[string]interface{}{
		"chat_id":    groupID,
		"text":       message,
		"parse_mode": "HTML",
	}
	
	if threadID != "" {
		payload["message_thread_id"] = threadID
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return
	}
	
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := t.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

// formatMessage creates formatted HTML message
func (t *TelegramNotifier) formatMessage(entry logger.LogEntry) string {
	var emoji string
	switch entry.Level {
	case logger.ERROR:
		emoji = "🚨"
	case logger.WARN:
		emoji = "⚠️"
	case logger.INFO:
		emoji = "ℹ️"
	default:
		emoji = "📝"
	}
	
	return fmt.Sprintf(
		"%s <b>%s</b>\n\n"+
		"<b>Time:</b> %s\n"+
		"<b>Context:</b> %s\n"+
		"<b>Function:</b> %s\n\n"+
		"<b>Message:</b>\n<code>%s</code>",
		emoji,
		entry.Level,
		entry.Timestamp.Format("2006-01-02 15:04:05"),
		entry.Context,
		entry.Function,
		entry.Message,
	)
}

// IsEnabled returns if notifier is configured and enabled
func (t *TelegramNotifier) IsEnabled() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.enabled
}
