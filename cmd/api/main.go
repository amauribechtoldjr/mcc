package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	handler "github.com/amauribechtoldjr/mcc/internal/adapter/handler/http"
	"github.com/amauribechtoldjr/mcc/internal/adapter/repository/postgres"
	repo "github.com/amauribechtoldjr/mcc/internal/adapter/repository/postgres/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/core/service"
	"github.com/amauribechtoldjr/mcc/internal/platform/env"
	"github.com/jackc/pgx/v5"
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

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		slog.Error("failed to connect to the database", "error", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	slog.Info("connected to database", "dsn", dsn)

	queries := repo.New(conn)

	cardRepo := postgres.NewCardRepository(queries)
	cardService := service.NewCardService(cardRepo)
	cardHandler := handler.NewCardHandler(cardService)

	collectionRepo := postgres.NewCollectionRepository(queries)
	collectionService := service.NewCollectionService(collectionRepo)
	collectionHandler := handler.NewCollectionHandler(collectionService)

	allowedOrigin := env.GetString("FE_ALLOWED_ORIGIN", "http://localhost:5173")

	router := handler.NewRouter(allowedOrigin, cardHandler, collectionHandler)

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
