package db

import (
	"database/sql"
	"fmt"

	"kafgres/internal/pkg/config"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// Interface abstracts database operations for testing.
type Interface interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Ping() error
	Close() error
}

// Connect opens a PostgreSQL connection.
func Connect(cfg config.PostgresConfig) (Interface, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	connection, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = connection.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logrus.Info("Successfully connected to PostgreSQL")
	return connection, nil
}
