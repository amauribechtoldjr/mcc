package main

import (
	"log"
	"net/http"
	"time"

	repo "github.com/amauribechtoldjr/mcc/internal/adapters/postgresql/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/domain/cards"
	"github.com/amauribechtoldjr/mcc/internal/domain/collections"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

type application struct {
	config config
	db     *pgx.Conn
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

func (a *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	// Cards
	cardsService := cards.NewService(repo.New(a.db))
	cardsHandler := cards.NewHandler(cardsService)

	r.Get("/cards", cardsHandler.ListCardsHandler)
	r.Get("/cards/{cardId}", cardsHandler.FindCardById)

	// Collections
	collectionsService := collections.NewService(repo.New(a.db))
	collectionsHandlers := collections.NewHandler(collectionsService)
	r.Post("/collections", collectionsHandlers.CreateCollection)

	return r
}

func (a *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         a.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at addr %s", a.config.addr)

	return srv.ListenAndServe()
}
