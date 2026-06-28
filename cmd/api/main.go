package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	repo "github.com/amauribechtoldjr/mcc/internal/db/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/handlers"
	"github.com/amauribechtoldjr/mcc/internal/repositories"
	"github.com/amauribechtoldjr/mcc/internal/sources"
	"github.com/amauribechtoldjr/mcc/internal/usecases"

	"github.com/amauribechtoldjr/mcc/internal/platform/database"
	"github.com/amauribechtoldjr/mcc/internal/platform/env"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		slog.Error("failed to start godotenv", "error", err)
		os.Exit(1)
	}

	dsn := env.GetString(
		"GOOSE_DBSTRING",
		"host=localhost user=postgres password=postgres dbname=mcc sslmode=disable",
	)

	conn, err := database.Connect(ctx, dsn)
	if err != nil {
		slog.Error("failed to connect to the database", "error", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	slog.Info("connected to database", "dsn", dsn)

	queries := repo.New(conn)

	repos := repositories.NewRepositories(queries)
	client := &http.Client{Timeout: 5 * time.Minute}
	source := sources.NewSource(client)
	useCases := usecases.New(repos, source)
	handlers := handlers.New(useCases)

	allowedOrigin := env.GetString("FE_ALLOWED_ORIGIN", "http://localhost:5173")

	router := handlers.NewRouter(allowedOrigin)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	slog.Info("server has started", "addr", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
