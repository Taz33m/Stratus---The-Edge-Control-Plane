package middleware

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LogEntry struct {
	RequestID  string        `json:"request_id"`
	Timestamp  time.Time     `json:"timestamp"`
	Method     string        `json:"method"`
	Path       string        `json:"path"`
	Status     int           `json:"status"`
	Latency    time.Duration `json:"latency_ms"`
	ClientIP   string        `json:"client_ip"`
	UserAgent  string        `json:"user_agent,omitempty"`
	Error      string        `json:"error,omitempty"`
}

func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		start := time.Now()

		// Process request
		c.Next()

		// Log after request
		entry := LogEntry{
			RequestID: requestID,
			Timestamp: start,
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Status:    c.Writer.Status(),
			Latency:   time.Since(start) / time.Millisecond,
			ClientIP:  c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		}

		if len(c.Errors) > 0 {
			entry.Error = c.Errors.String()
		}

		// Output as JSON
		logJSON, _ := json.Marshal(entry)
		println(string(logJSON))
	}
}
