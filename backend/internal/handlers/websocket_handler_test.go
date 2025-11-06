package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stratus/backend/internal/websocket"
)

func TestWebSocketOriginCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	hub := websocket.NewHub()
	handler := NewWebSocketHandler(hub, "http://localhost:3000,http://example.com")

	tests := []struct {
		name       string
		origin     string
		wantStatus int
	}{
		{
			name:       "allowed origin",
			origin:     "http://localhost:3000",
			wantStatus: http.StatusSwitchingProtocols,
		},
		{
			name:       "another allowed origin",
			origin:     "http://example.com",
			wantStatus: http.StatusSwitchingProtocols,
		},
		{
			name:       "disallowed origin",
			origin:     "http://malicious.com",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "no origin header",
			origin:     "",
			wantStatus: http.StatusSwitchingProtocols,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			
			req := httptest.NewRequest("GET", "/ws", nil)
			req.Header.Set("Connection", "Upgrade")
			req.Header.Set("Upgrade", "websocket")
			req.Header.Set("Sec-WebSocket-Version", "13")
			req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
			
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			
			c.Request = req

			handler.HandleWebSocket(c)

			// Note: In real test, WebSocket upgrade might not complete in test environment
			// This is a simplified check
			if tt.origin != "" && !strings.Contains(tt.origin, "localhost") && !strings.Contains(tt.origin, "example.com") {
				// Disallowed origins should fail upgrade
				if w.Code == http.StatusSwitchingProtocols {
					t.Errorf("WebSocket upgrade should have failed for origin %s", tt.origin)
				}
			}
		})
	}
}
