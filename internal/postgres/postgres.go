package postgres

import (
	"context"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/DeluxeOwl/cogniboard/internal/postgres/ent"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

// Very primitive support for migrations for the quick demo - don't use in production
func NewPostgresWithMigrate(ctx context.Context, dsn string) (*ent.Client, error) {
	db, err := connectDB(dsn)
	if err != nil {
		return nil, err
	}

	client := ent.NewClient(ent.Driver(db))
	if err := client.Schema.Create(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

func connectDB(connString string) (*sql.Driver, error) {
	config, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	db := stdlib.OpenDB(*config)

	// Create an ent driver
	drv := sql.OpenDB(dialect.Postgres, db)

	return drv, nil
}
