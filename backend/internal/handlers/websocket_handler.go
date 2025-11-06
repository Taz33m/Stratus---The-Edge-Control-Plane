package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	ws "github.com/stratus/backend/internal/websocket"
)

type WebSocketHandler struct {
	hub            *ws.Hub
	allowedOrigins map[string]bool
}

func NewWebSocketHandler(hub *ws.Hub, corsOrigins string) *WebSocketHandler {
	allowedOrigins := make(map[string]bool)
	for _, origin := range strings.Split(corsOrigins, ",") {
		allowedOrigins[strings.TrimSpace(origin)] = true
	}
	return &WebSocketHandler{
		hub:            hub,
		allowedOrigins: allowedOrigins,
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return true // Allow same-origin requests
			}
			return h.allowedOrigins[origin]
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := ws.NewClient(h.hub, conn)
	h.hub.Register(client)
	
	go client.WritePump()
	go client.ReadPump()
}
