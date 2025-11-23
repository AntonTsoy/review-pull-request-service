package database

import (
	"context"
	"fmt"
	"time"

	"github.com/AntonTsoy/review-pull-request-service/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, dbConfig *config.Config) (*Database, error) {
	poolConfig, err := pgxpool.ParseConfig(connectionString(dbConfig))
	if err != nil {
		return nil, fmt.Errorf("parse database connection config: %w", err)
	}

	poolConfig.MaxConns = dbConfig.Pool

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	return &Database{pool: pool}, nil
}

func (db *Database) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *Database) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return db.pool.BeginTx(ctx, txOptions)
}

func (db *Database) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return db.pool.Ping(ctx)
}

func (db *Database) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

func connectionString(dbConfig *config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		dbConfig.SSL,
	)
}
