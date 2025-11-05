package models

import (
	"time"
)

type ServiceStatus string

const (
	StatusRunning ServiceStatus = "running"
	StatusStopped ServiceStatus = "stopped"
	StatusError   ServiceStatus = "error"
	StatusStarting ServiceStatus = "starting"
)

type Service struct {
	ID        string        `json:"id" db:"id"`
	Name      string        `json:"name" db:"name"`
	Region    string        `json:"region" db:"region"`
	Image     string        `json:"image" db:"image"`
	Version   string        `json:"version" db:"version"`
	Status    ServiceStatus `json:"status" db:"status"`
	Uptime    int64         `json:"uptime" db:"uptime"` // seconds
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
}

type CreateServiceRequest struct {
	Name    string `json:"name" binding:"required"`
	Region  string `json:"region" binding:"required"`
	Image   string `json:"image" binding:"required"`
	Version string `json:"version" binding:"required"`
}

type UpdateServiceRequest struct {
	Status  *ServiceStatus `json:"status,omitempty"`
	Version *string        `json:"version,omitempty"`
}

type ServiceMetrics struct {
	ServiceID    string    `json:"service_id"`
	Timestamp    time.Time `json:"timestamp"`
	CPUUsage     float64   `json:"cpu_usage"`      // percentage
	MemoryUsage  float64   `json:"memory_usage"`   // MB
	RequestCount int64     `json:"request_count"`  // total requests
	ErrorRate    float64   `json:"error_rate"`     // percentage
	P95Latency   float64   `json:"p95_latency"`    // milliseconds (p95)
}

type ServiceConfig struct {
	ID          string                 `json:"id" db:"id"`
	ServiceID   string                 `json:"service_id" db:"service_id"`
	Config      map[string]interface{} `json:"config" db:"config"`
	Version     int                    `json:"version" db:"version"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	CreatedBy   string                 `json:"created_by" db:"created_by"`
}

type DeploymentLog struct {
	ID        string    `json:"id" db:"id"`
	ServiceID string    `json:"service_id" db:"service_id"`
	Action    string    `json:"action" db:"action"` // start, stop, restart, config_update
	Status    string    `json:"status" db:"status"` // pending, success, failed
	Message   string    `json:"message" db:"message"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
