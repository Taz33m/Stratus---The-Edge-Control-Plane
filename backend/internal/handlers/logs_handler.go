package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stratus/backend/internal/models"
)

type LogsHandler struct {
	db *sql.DB
}

func NewLogsHandler(db *sql.DB) *LogsHandler {
	return &LogsHandler{db: db}
}

func (h *LogsHandler) GetDeploymentLogs(c *gin.Context) {
	limit := 50
	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	query := `
		SELECT dl.id, dl.service_id, s.name as service_name, dl.action, dl.status, dl.message, dl.created_at
		FROM deployment_logs dl
		LEFT JOIN services s ON dl.service_id = s.id
		ORDER BY dl.created_at DESC
		LIMIT $1
	`

	rows, err := h.db.Query(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type DeploymentLogWithName struct {
		models.DeploymentLog
		ServiceName string `json:"service_name"`
	}

	logs := []DeploymentLogWithName{}
	for rows.Next() {
		var log DeploymentLogWithName
		if err := rows.Scan(&log.ID, &log.ServiceID, &log.ServiceName, &log.Action, &log.Status, &log.Message, &log.CreatedAt); err != nil {
			continue
		}
		logs = append(logs, log)
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
