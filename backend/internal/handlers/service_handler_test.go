package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stratus/backend/internal/models"
	"github.com/stratus/backend/internal/websocket"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Use in-memory SQLite for testing or mock
	// For now, return nil and skip DB-dependent tests
	return nil
}

func TestCreateServiceValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	hub := websocket.NewHub()
	metricsHandler := &MetricsHandler{
		simulators: make(map[string]context.CancelFunc),
	}
	handler := NewServiceHandler(nil, hub, metricsHandler)

	tests := []struct {
		name       string
		payload    models.CreateServiceRequest
		wantStatus int
	}{
		{
			name: "valid service",
			payload: models.CreateServiceRequest{
				Name:    "test-service",
				Region:  "us-east-1",
				Image:   "nginx",
				Version: "1.0.0",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid region",
			payload: models.CreateServiceRequest{
				Name:    "test-service",
				Region:  "invalid-region",
				Image:   "nginx",
				Version: "1.0.0",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid service name",
			payload: models.CreateServiceRequest{
				Name:    "a", // too short
				Region:  "us-east-1",
				Image:   "nginx",
				Version: "1.0.0",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if handler.db == nil {
				t.Skip("Skipping test that requires database")
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.payload)
			c.Request = httptest.NewRequest("POST", "/api/v1/services", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.CreateService(c)

			if w.Code != tt.wantStatus {
				t.Errorf("CreateService() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestListServicesPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	hub := websocket.NewHub()
	metricsHandler := &MetricsHandler{
		simulators: make(map[string]context.CancelFunc),
	}
	handler := NewServiceHandler(nil, hub, metricsHandler)

	tests := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{
			name:       "default pagination",
			query:      "",
			wantStatus: http.StatusOK,
		},
		{
			name:       "with limit",
			query:      "?limit=10",
			wantStatus: http.StatusOK,
		},
		{
			name:       "with offset",
			query:      "?offset=20",
			wantStatus: http.StatusOK,
		},
		{
			name:       "with limit and offset",
			query:      "?limit=25&offset=50",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if handler.db == nil {
				t.Skip("Skipping test that requires database")
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/api/v1/services"+tt.query, nil)

			handler.ListServices(c)

			if w.Code != tt.wantStatus {
				t.Errorf("ListServices() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}
