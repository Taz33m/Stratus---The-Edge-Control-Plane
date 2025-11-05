package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stratus/backend/internal/models"
	"github.com/stratus/backend/internal/websocket"
)

type ServiceHandler struct {
	db  *sql.DB
	hub *websocket.Hub
}

func NewServiceHandler(db *sql.DB, hub *websocket.Hub) *ServiceHandler {
	return &ServiceHandler{
		db:  db,
		hub: hub,
	}
}

func (h *ServiceHandler) ListServices(c *gin.Context) {
	region := c.Query("region")
	status := c.Query("status")

	query := "SELECT id, name, region, image, version, status, uptime, created_at, updated_at FROM services WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if region != "" {
		query += " AND region = $" + string(rune(argCount+'0'))
		args = append(args, region)
		argCount++
	}

	if status != "" {
		query += " AND status = $" + string(rune(argCount+'0'))
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	services := []models.Service{}
	for rows.Next() {
		var s models.Service
		if err := rows.Scan(&s.ID, &s.Name, &s.Region, &s.Image, &s.Version, &s.Status, &s.Uptime, &s.CreatedAt, &s.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		services = append(services, s)
	}

	c.JSON(http.StatusOK, gin.H{"services": services})
}

func (h *ServiceHandler) GetService(c *gin.Context) {
	id := c.Param("id")

	var s models.Service
	err := h.db.QueryRow(
		"SELECT id, name, region, image, version, status, uptime, created_at, updated_at FROM services WHERE id = $1",
		id,
	).Scan(&s.ID, &s.Name, &s.Region, &s.Image, &s.Version, &s.Status, &s.Uptime, &s.CreatedAt, &s.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s)
}

func (h *ServiceHandler) CreateService(c *gin.Context) {
	var req models.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service := models.Service{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Region:    req.Region,
		Image:     req.Image,
		Version:   req.Version,
		Status:    models.StatusStopped,
		Uptime:    0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := h.db.Exec(
		`INSERT INTO services (id, name, region, image, version, status, uptime, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		service.ID, service.Name, service.Region, service.Image, service.Version, service.Status, service.Uptime, service.CreatedAt, service.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Broadcast service creation
	h.hub.BroadcastJSON(websocket.MessageTypeServiceUpdate, service)

	// Create deployment log
	h.createDeploymentLog(service.ID, "create", "success", "Service created successfully")

	c.JSON(http.StatusCreated, service)
}

func (h *ServiceHandler) UpdateService(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build update query dynamically
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if req.Status != nil {
		updates = append(updates, "status = $"+string(rune(argCount+'0')))
		args = append(args, *req.Status)
		argCount++
	}

	if req.Version != nil {
		updates = append(updates, "version = $"+string(rune(argCount+'0')))
		args = append(args, *req.Version)
		argCount++
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	updates = append(updates, "updated_at = $"+string(rune(argCount+'0')))
	args = append(args, time.Now())
	argCount++

	args = append(args, id)

	query := "UPDATE services SET "
	for i, update := range updates {
		if i > 0 {
			query += ", "
		}
		query += update
	}
	query += " WHERE id = $" + string(rune(argCount+'0'))

	result, err := h.db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	// Fetch updated service
	var service models.Service
	h.db.QueryRow(
		"SELECT id, name, region, image, version, status, uptime, created_at, updated_at FROM services WHERE id = $1",
		id,
	).Scan(&service.ID, &service.Name, &service.Region, &service.Image, &service.Version, &service.Status, &service.Uptime, &service.CreatedAt, &service.UpdatedAt)

	// Broadcast update
	h.hub.BroadcastJSON(websocket.MessageTypeServiceUpdate, service)

	action := "update"
	if req.Status != nil {
		switch *req.Status {
		case models.StatusRunning:
			action = "start"
		case models.StatusStopped:
			action = "stop"
		}
	}

	h.createDeploymentLog(id, action, "success", "Service updated successfully")

	c.JSON(http.StatusOK, service)
}

func (h *ServiceHandler) DeleteService(c *gin.Context) {
	id := c.Param("id")

	result, err := h.db.Exec("DELETE FROM services WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	h.hub.Broadcast(websocket.MessageTypeServiceUpdate, gin.H{
		"id":     id,
		"action": "deleted",
	})

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}

func (h *ServiceHandler) createDeploymentLog(serviceID, action, status, message string) {
	logID := uuid.New().String()
	h.db.Exec(
		`INSERT INTO deployment_logs (id, service_id, action, status, message, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		logID, serviceID, action, status, message, time.Now(),
	)

	log := models.DeploymentLog{
		ID:        logID,
		ServiceID: serviceID,
		Action:    action,
		Status:    status,
		Message:   message,
		CreatedAt: time.Now(),
	}

	h.hub.BroadcastJSON(websocket.MessageTypeLog, log)
}
