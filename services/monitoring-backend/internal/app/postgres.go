package app

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/qreaqtor/containers-monitoring/monitoring-backend/internal/config"
)

func getPostgresConn(cfg config.PostgresConfig) (*sql.DB, error) {
	conn, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to PostgreSQL: %v", err)
	}

	return conn, nil
}
