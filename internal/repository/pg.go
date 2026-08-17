package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Database(ctx context.Context, url string) (*pgxpool.Pool, error) { return pgxpool.New(ctx, url) }
