package handlers

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stratus/backend/internal/models"
	"github.com/stratus/backend/internal/websocket"
)

type MetricsHandler struct {
	redis      *redis.Client
	hub        *websocket.Hub
	simulators map[string]context.CancelFunc
	mu         sync.RWMutex
}

func NewMetricsHandler(redisClient *redis.Client, hub *websocket.Hub) *MetricsHandler {
	return &MetricsHandler{
		redis:      redisClient,
		hub:        hub,
		simulators: make(map[string]context.CancelFunc),
	}
}

func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	serviceID := c.Param("id")
	ctx := context.Background()

	// Get latest metrics from Redis
	key := "metrics:" + serviceID
	data, err := h.redis.LRange(ctx, key, 0, 99).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metrics := []models.ServiceMetrics{}
	for _, item := range data {
		var m models.ServiceMetrics
		if err := json.Unmarshal([]byte(item), &m); err != nil {
			continue
		}
		metrics = append(metrics, m)
	}

	c.JSON(http.StatusOK, gin.H{"metrics": metrics})
}

// StartSimulator starts metrics simulation for a service
func (h *MetricsHandler) StartSimulator(serviceID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Stop existing simulator if any
	if cancel, exists := h.simulators[serviceID]; exists {
		cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	h.simulators[serviceID] = cancel

	go h.simulateMetrics(ctx, serviceID)
}

// StopSimulator stops metrics simulation for a service
func (h *MetricsHandler) StopSimulator(serviceID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if cancel, exists := h.simulators[serviceID]; exists {
		cancel()
		delete(h.simulators, serviceID)
	}
}

func (h *MetricsHandler) simulateMetrics(ctx context.Context, serviceID string) {
	key := "metrics:" + serviceID
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			metrics := models.ServiceMetrics{
				ServiceID:    serviceID,
				Timestamp:    time.Now(),
				CPUUsage:     rand.Float64() * 80,          // 0-80%
				MemoryUsage:  rand.Float64()*500 + 100,     // 100-600 MB (actual MB, not %)
				RequestCount: rand.Int63n(1000),            // Total requests in this interval
				ErrorRate:    rand.Float64() * 5,           // 0-5%
				P95Latency:   rand.Float64()*200 + 50,      // 50-250ms
			}

			data, _ := json.Marshal(metrics)
			h.redis.LPush(ctx, key, data)
			h.redis.LTrim(ctx, key, 0, 99) // Keep last 100 metrics
			h.redis.Expire(ctx, key, 24*time.Hour)

			// Broadcast metrics
			h.hub.BroadcastJSON(websocket.MessageTypeMetrics, metrics)
		}
	}
}

func (h *MetricsHandler) GetAggregatedMetrics(c *gin.Context) {
	// This would aggregate metrics across all services
	// For now, return mock aggregated data
	
	aggregated := gin.H{
		"total_services": 12,
		"running_services": 10,
		"avg_cpu_usage": 45.3,
		"avg_memory_usage": 320.5,
		"total_requests": 125000,
		"avg_error_rate": 0.8,
		"regions": gin.H{
			"US-East": 4,
			"US-West": 3,
			"EU-West": 3,
			"APAC": 2,
		},
	}

	c.JSON(http.StatusOK, aggregated)
}
