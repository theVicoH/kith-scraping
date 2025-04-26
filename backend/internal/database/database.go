package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"backend/internal/config"
	"backend/internal/config/types"

	_ "github.com/lib/pq"
)

type Connection struct {
	DB *sql.DB
}

func NewDatabaseClient(cfg *config.Config) (*Connection, error) {
	db, err := sql.Open("postgres", dns(cfg.Database))
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database connection ping failed: %w", err)
	}

	fmt.Println("Database connection established successfully")

	return &Connection{DB: db}, nil
}

func dns(dc types.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dc.Host, dc.Port, dc.User, dc.Password, dc.Name,
	)
}
