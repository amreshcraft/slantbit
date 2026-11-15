package database

import (
	"context"
	"fmt"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgres(dbURL string) (*pgxpool.Pool, error) {
	// Parse connection string
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	//  Production Pool Settings
	config.MaxConns = 25                      // Ek saath maximum 25 connections
	config.MinConns = 5                       // Hamesha 5 connections ready rahenge
	config.MaxConnLifetime = time.Hour        // 1 ghante baad connection refresh hoga
	config.MaxConnIdleTime = 30 * time.Minute // 30 min khali rehne par close

	// 3. Connection timeout ke saath pool banayein
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// 4. Verify connection (Ping)
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return pool, nil
}