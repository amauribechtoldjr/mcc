package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/amauribechtoldjr/mcc/internal/env"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		slog.Error("failed to start godotenv", "error", err)
		os.Exit(1)
	}

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString(
				"GOOSE_DBSTRING",
				"host=localhost user=postgres password=postgres dbname=mcc sslmode=disable",
			),
		},
	}

	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		slog.Error("failed to connect to the database", "error", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	slog.Info("connected to database", "dsn", cfg.db.dsn)

	api := application{
		config: cfg,
		db:     conn,
	}

	err = api.run(api.mount())
	if err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
