package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"user-service/internal/entity"
)

type Storage struct {
	pool *pgxpool.Pool
	log  entity.Logger
}

func New(connStr string, log entity.Logger) (*Storage, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return &Storage{
		pool: pool,
		log:  log,
	}, nil
}
