package db

import (
	"context"
	"sport_platform/internal/sqlc/db_queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbClient struct {
	connectionPool *pgxpool.Pool
	Queries        *db_queries.Queries
}

func (dbClient DbClient) Close() error {
	dbClient.connectionPool.Close()
	return nil
}

func (dbClient DbClient) Exec(ctx context.Context, sql string, args ...interface{}) error {
	_, err := dbClient.connectionPool.Exec(ctx, sql, args...)
	return err
}

func CreateConnection(cfg *Config, ctx context.Context) (*DbClient, error) {
	connectionPool, poolCreationError := pgxpool.New(ctx, cfg.GetConnectionString())
	if poolCreationError != nil {
		return nil, poolCreationError
	}

	return &DbClient{
		connectionPool: connectionPool,
		Queries:        db_queries.New(connectionPool),
	}, nil
}
