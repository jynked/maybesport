package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:4r5t_6Y7U@localhost:5432/maybesport?sslmode=disable"
	}
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return err
	}
	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	return err
}
