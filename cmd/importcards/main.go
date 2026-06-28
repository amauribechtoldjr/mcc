package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	repo "github.com/amauribechtoldjr/mcc/internal/db/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/platform/database"
	"github.com/amauribechtoldjr/mcc/internal/platform/env"
	"github.com/amauribechtoldjr/mcc/internal/repositories"
	"github.com/amauribechtoldjr/mcc/internal/sources"
	"github.com/amauribechtoldjr/mcc/internal/usecases"
	"github.com/joho/godotenv"
)

const (
	gameCode  = "mtg"
	cardLimit = 900000
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

	repos := repositories.NewRepositories(queries)

	client := &http.Client{Timeout: 5 * time.Minute}
	source := sources.NewSource(client)

	usecases := usecases.New(repos, source)

	if err := usecases.ScryfallImport.Run(ctx, gameCode, cardLimit); err != nil {
		slog.Error("failed to import cards", "error", err)
		os.Exit(1)
	}

	slog.Info("import command finished")
}
