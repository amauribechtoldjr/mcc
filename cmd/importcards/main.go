package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/amauribechtoldjr/mcc/internal/adapter/repository/postgres"
	repo "github.com/amauribechtoldjr/mcc/internal/adapter/repository/postgres/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/adapter/scryfall"
	"github.com/amauribechtoldjr/mcc/internal/core/service"
	"github.com/amauribechtoldjr/mcc/internal/platform/database"
	"github.com/amauribechtoldjr/mcc/internal/platform/env"
	"github.com/joho/godotenv"
)

const (
	gameCode  = "mtg"
	cardLimit = 900000
	userAgent = "my-magic-collection/0.1 (contact: amauribechtoldjr@gmail.com)"
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

	slog.Info("connected to database")

	queries := repo.New(conn)

	cardRepo := postgres.NewCardRepository(queries)
	gameRepo := postgres.NewGameRepository(queries)
	mtgSetRepo := postgres.NewMTGSetRepository(queries)
	scryfallImportRepo := postgres.NewScryfallImportRepository(queries)

	client := &http.Client{Timeout: 5 * time.Minute}
	source := scryfall.NewCardSource(client, userAgent)

	importService := service.NewImportService(source, cardRepo, gameRepo, mtgSetRepo, scryfallImportRepo)

	if err := importService.Run(ctx, gameCode, cardLimit); err != nil {
		slog.Error("failed to import cards", "error", err)
		os.Exit(1)
	}

	slog.Info("import command finished")
}
