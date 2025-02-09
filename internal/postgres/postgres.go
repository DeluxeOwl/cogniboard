package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var migrations = []string{`
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    description TEXT,
    due_date TIMESTAMP WITH TIME ZONE,
    assignee TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) NOT NULL
);
`, `
CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    mime_type VARCHAR(127) NOT NULL,
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL
);
`, `
CREATE TABLE IF NOT EXISTS task_files (
    task_id UUID REFERENCES tasks(id) ON DELETE CASCADE,
    file_id UUID REFERENCES files(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, file_id)
);
`}

// Very primitive support for migrations for the quick demo - don't use in production
func NewPostgresWithMigrate(dsn string) (*sqlx.DB, error) {
	db, err := connectDB(dsn)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	for _, migration := range migrations {
		tx.ExecContext(ctx, migration)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return db, nil
}

// connectDB creates a new sqlx.DB instance with connection settings
func connectDB(dsn string) (*sqlx.DB, error) {
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse DSN: %w", err)
	}

	// Set timeouts
	config.ConnectTimeout = 10 * time.Second
	config.RuntimeParams["statement_timeout"] = "30000" // 30 seconds
	config.RuntimeParams["idle_in_transaction_session_timeout"] = "30000"

	connStr := stdlib.RegisterConnConfig(config)
	db, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}
