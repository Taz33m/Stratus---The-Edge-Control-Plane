package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stratus/backend/internal/errors"
	"github.com/stratus/backend/internal/models"
	"github.com/stratus/backend/internal/validation"
	"github.com/stratus/backend/internal/websocket"
)

type ServiceHandler struct {
	db      *sql.DB
	hub     *websocket.Hub
	metrics *MetricsHandler
}

func NewServiceHandler(db *sql.DB, hub *websocket.Hub, metrics *MetricsHandler) *ServiceHandler {
	return &ServiceHandler{
		db:      db,
		hub:     hub,
		metrics: metrics,
	}
}

func (h *ServiceHandler) ListServices(c *gin.Context) {
	region := c.Query("region")
	status := c.Query("status")

	// Pagination
	limit := 50
	offset := 0
	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := fmt.Sscanf(limitParam, "%d", &limit); err == nil && l == 1 && limit > 0 && limit <= 100 {
			// Valid limit
		} else {
			limit = 50
		}
	}
	if offsetParam := c.Query("offset"); offsetParam != "" {
		if o, err := fmt.Sscanf(offsetParam, "%d", &offset); err == nil && o == 1 && offset >= 0 {
			// Valid offset
		} else {
			offset = 0
		}
	}

	query := "SELECT id, name, region, image, version, status, uptime, created_at, updated_at FROM services WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if region != "" {
		query += fmt.Sprintf(" AND region = $%d", argCount)
		args = append(args, region)
		argCount++
	}

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, query, args...)
	if err != nil {
		errors.InternalError(c, "Failed to query services")
		return
	}
	defer rows.Close()

	services := []models.Service{}
	for rows.Next() {
		var s models.Service
		if err := rows.Scan(&s.ID, &s.Name, &s.Region, &s.Image, &s.Version, &s.Status, &s.Uptime, &s.CreatedAt, &s.UpdatedAt); err != nil {
			errors.InternalError(c, "Failed to scan service")
			return
		}
		services = append(services, s)
	}

	c.JSON(http.StatusOK, gin.H{"services": services})
}

func (h *ServiceHandler) GetService(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var s models.Service
	err := h.db.QueryRowContext(ctx,
		"SELECT id, name, region, image, version, status, uptime, created_at, updated_at FROM services WHERE id = $1",
		id,
	).Scan(&s.ID, &s.Name, &s.Region, &s.Image, &s.Version, &s.Status, &s.Uptime, &s.CreatedAt, &s.UpdatedAt)

	if err == sql.ErrNoRows {
		errors.NotFound(c, "Service")
		return
	}
	if err != nil {
		errors.InternalError(c, "Failed to get service")
		return
	}

	c.JSON(http.StatusOK, s)
}

func (h *ServiceHandler) CreateService(c *gin.Context) {
	var req models.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Validate input
	var validationErrs validation.ValidationErrors
	if err := validation.ValidateServiceName(req.Name); err != nil {
		validationErrs = append(validationErrs, err.(validation.ValidationError))
	}
	if err := validation.ValidateRegion(req.Region); err != nil {
		validationErrs = append(validationErrs, err.(validation.ValidationError))
	}
	if err := validation.ValidateImage(req.Image); err != nil {
		validationErrs = append(validationErrs, err.(validation.ValidationError))
	}
	if err := validation.ValidateVersion(req.Version); err != nil {
		validationErrs = append(validationErrs, err.(validation.ValidationError))
	}
	if len(validationErrs) > 0 {
		errors.BadRequest(c, "Validation failed", validationErrs)
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	_, err := h.db.ExecContext(ctx,
		`INSERT INTO services (id, name, region, image, version, status, uptime, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		service.ID, service.Name, service.Region, service.Image, service.Version, service.Status, service.Uptime, service.CreatedAt, service.UpdatedAt,
	)

	if err != nil {
		errors.InternalError(c, "Failed to create service")
		return
	}

	// Broadcast service creation
	h.hub.BroadcastJSON(websocket.MessageTypeServiceUpdate, service)

	// Create deployment log
	h.createDeploymentLog(service.ID, "create", "success", "Service created successfully")

	// Start metrics simulator for new service
	h.metrics.StartSimulator(service.ID)

	c.JSON(http.StatusCreated, service)
}

func (h *ServiceHandler) UpdateService(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Validate version if provided
	if req.Version != nil {
		if err := validation.ValidateVersion(*req.Version); err != nil {
			errors.BadRequest(c, "Validation failed", err)
			return
		}
	}

	// Build update query dynamically
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if req.Status != nil {
		updates = append(updates, fmt.Sprintf("status = $%d", argCount))
		args = append(args, *req.Status)
		argCount++
	}

	if req.Version != nil {
		updates = append(updates, fmt.Sprintf("version = $%d", argCount))
		args = append(args, *req.Version)
		argCount++
	}

	if len(updates) == 0 {
		errors.BadRequest(c, "No fields to update", nil)
		return
	}

	updates = append(updates, fmt.Sprintf("updated_at = $%d", argCount))
	args = append(args, time.Now())
	argCount++

	args = append(args, id)

	query := "UPDATE services SET " + strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d", argCount)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	result, err := h.db.ExecContext(ctx, query, args...)
	if err != nil {
		errors.InternalError(c, "Failed to update service")
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		errors.NotFound(c, "Service")
		return
	}

	// Fetch updated service
	var service models.Service
	h.db.QueryRowContext(ctx,
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
			// Start metrics simulator when service starts
			h.metrics.StartSimulator(id)
		case models.StatusStopped:
			action = "stop"
			// Stop metrics simulator when service stops
			h.metrics.StopSimulator(id)
		}
	}

	h.createDeploymentLog(id, action, "success", "Service updated successfully")

	c.JSON(http.StatusOK, service)
}

func (h *ServiceHandler) DeleteService(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	result, err := h.db.ExecContext(ctx, "DELETE FROM services WHERE id = $1", id)
	if err != nil {
		errors.InternalError(c, "Failed to delete service")
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		errors.NotFound(c, "Service")
		return
	}

	// Stop metrics simulator
	h.metrics.StopSimulator(id)

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
