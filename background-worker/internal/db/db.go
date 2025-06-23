package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect(ctx context.Context, connStr string) error {
	var err error
	Pool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		return err
	}

	if err := Pool.Ping(ctx); err != nil {
		return err
	}

	return nil
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}

