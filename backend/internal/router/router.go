package router

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stratus/backend/internal/config"
	"github.com/stratus/backend/internal/handlers"
	"github.com/stratus/backend/internal/middleware"
	"github.com/stratus/backend/internal/websocket"
)

func Setup(cfg *config.Config, db *sql.DB, redisClient *redis.Client, hub *websocket.Hub) *gin.Engine {
	r := gin.New()

	// Recovery middleware
	r.Use(gin.Recovery())

	// Structured logging
	r.Use(middleware.StructuredLogger())

	// Secure CORS middleware
	r.Use(func(c *gin.Context) {
		origins := strings.Split(cfg.CORSOrigins, ",")
		origin := c.Request.Header.Get("Origin")
		allowed := false
		
		for _, allowedOrigin := range origins {
			allowedOrigin = strings.TrimSpace(allowedOrigin)
			if origin == allowedOrigin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				allowed = true
				break
			}
		}
		
		if !allowed && origin != "" {
			// Don't set CORS headers for disallowed origins
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(403)
				return
			}
		}
		
		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize middleware
	auth := middleware.NewAuthMiddleware(cfg.JWTSecret)
	rateLimiter := middleware.NewRateLimiter(100, time.Minute) // 100 requests per minute

	// Initialize handlers
	metricsHandler := handlers.NewMetricsHandler(redisClient, hub)
	serviceHandler := handlers.NewServiceHandler(db, hub, metricsHandler)
	wsHandler := handlers.NewWebSocketHandler(hub, cfg.CORSOrigins)
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
		// Public/viewer endpoints (read-only)
		public := v1.Group("")
		public.Use(auth.AuthRequired())
		public.Use(auth.RequireRole(middleware.RoleViewer))
		{
			public.GET("/services", serviceHandler.ListServices)
			public.GET("/services/:id", serviceHandler.GetService)
			public.GET("/metrics/:id", metricsHandler.GetMetrics)
			public.GET("/metrics/aggregated", metricsHandler.GetAggregatedMetrics)
			public.GET("/logs/deployment", logsHandler.GetDeploymentLogs)
		}

		// Operator endpoints (mutating operations)
		operator := v1.Group("")
		operator.Use(auth.AuthRequired())
		operator.Use(auth.RequireRole(middleware.RoleOperator))
		operator.Use(rateLimiter.Limit()) // Rate limit mutating operations
		{
			operator.POST("/services", serviceHandler.CreateService)
			operator.PATCH("/services/:id", serviceHandler.UpdateService)
		}

		// Admin endpoints
		admin := v1.Group("")
		admin.Use(auth.AuthRequired())
		admin.Use(auth.RequireRole(middleware.RoleAdmin))
		admin.Use(rateLimiter.Limit())
		{
			admin.DELETE("/services/:id", serviceHandler.DeleteService)
		}
	}

	// WebSocket endpoint (requires auth)
	ws := r.Group("")
	ws.Use(auth.AuthRequired())
	{
		ws.GET("/ws", wsHandler.HandleWebSocket)
	}

	// Token generation endpoint (for development/testing)
	if cfg.Environment == "development" {
		r.POST("/auth/token", func(c *gin.Context) {
			var req struct {
				UserID string `json:"user_id" binding:"required"`
				Role   string `json:"role" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			token, err := auth.GenerateToken(req.UserID, middleware.Role(req.Role))
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"token": token})
		})
	}

	return r
}
