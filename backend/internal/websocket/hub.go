package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

type MessageType string

const (
	MessageTypeServiceUpdate MessageType = "service_update"
	MessageTypeMetrics       MessageType = "metrics"
	MessageTypeLog           MessageType = "log"
)

type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Broadcast(msgType MessageType, payload interface{}) {
	message := Message{
		Type:    msgType,
		Payload: payload,
	}
	h.broadcast <- message
}

func (h *Hub) BroadcastJSON(msgType MessageType, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	
	var jsonPayload interface{}
	if err := json.Unmarshal(data, &jsonPayload); err != nil {
		return err
	}
	
	h.Broadcast(msgType, jsonPayload)
	return nil
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}
