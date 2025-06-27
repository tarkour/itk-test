package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tarkour/itk-test/pkg/config"
)

func ConnectDB(path string) (*DBConn, error) {

	cfg, err := config.LoadConfig(path)
	if err != nil {
		log.Fatal("Config error: ", err)
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.DBSource)
	if err != nil {
		return nil, fmt.Errorf("pool config error: %w", err)
	}

	poolConfig.MinConns = 10
	poolConfig.MaxConns = 100
	poolConfig.MaxConnLifetime = 0
	poolConfig.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("connection pool error: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("database test ping error: %v", err)
	}

	return &DBConn{Pool: pool}, nil
}

func (db *DBConn) Close() {
	db.Pool.Close()
}
