package router

import (
	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stratus/backend/internal/config"
	"github.com/stratus/backend/internal/handlers"
	"github.com/stratus/backend/internal/websocket"
)

func Setup(cfg *config.Config, db *sql.DB, redisClient *redis.Client, hub *websocket.Hub) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		origins := strings.Split(cfg.CORSOrigins, ",")
		origin := c.Request.Header.Get("Origin")
		
		for _, allowedOrigin := range origins {
			if origin == strings.TrimSpace(allowedOrigin) || allowedOrigin == "*" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
		
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize handlers
	serviceHandler := handlers.NewServiceHandler(db, hub)
	metricsHandler := handlers.NewMetricsHandler(redisClient, hub)
	wsHandler := handlers.NewWebSocketHandler(hub)
	logsHandler := handlers.NewLogsHandler(db)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "stratus-control-plane",
		})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Services
		services := v1.Group("/services")
		{
			services.GET("", serviceHandler.ListServices)
			services.POST("", serviceHandler.CreateService)
			services.GET("/:id", serviceHandler.GetService)
			services.PATCH("/:id", serviceHandler.UpdateService)
			services.DELETE("/:id", serviceHandler.DeleteService)
		}

		// Metrics
		metrics := v1.Group("/metrics")
		{
			metrics.GET("/:id", metricsHandler.GetMetrics)
			metrics.GET("/aggregated", metricsHandler.GetAggregatedMetrics)
		}

		// Logs
		logs := v1.Group("/logs")
		{
			logs.GET("/deployment", logsHandler.GetDeploymentLogs)
		}
	}

	// WebSocket endpoint
	r.GET("/ws", wsHandler.HandleWebSocket)

	return r
}
