package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPostgresDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS services (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			region VARCHAR(50) NOT NULL,
			image VARCHAR(255) NOT NULL,
			version VARCHAR(50) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'stopped',
			uptime BIGINT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS service_configs (
			id VARCHAR(36) PRIMARY KEY,
			service_id VARCHAR(36) NOT NULL,
			config JSONB NOT NULL,
			version INT NOT NULL DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(255),
			FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS deployment_logs (
			id VARCHAR(36) PRIMARY KEY,
			service_id VARCHAR(36) NOT NULL,
			action VARCHAR(50) NOT NULL,
			status VARCHAR(20) NOT NULL,
			message TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_services_region ON services(region)`,
		`CREATE INDEX IF NOT EXISTS idx_services_status ON services(status)`,
		`CREATE INDEX IF NOT EXISTS idx_deployment_logs_service_id ON deployment_logs(service_id)`,
		`CREATE INDEX IF NOT EXISTS idx_deployment_logs_created_at ON deployment_logs(created_at DESC)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}
