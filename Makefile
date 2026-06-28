dev:
	air

migrate-add:
	goose -s create $(name) sql -dir "./internal/db/migrations"

migrate-up:
	goose up

migrate-down:
	goose down

migrate-down-to:
	goose down-to $(version)

sql-gen:
	sqlc generate

start-infra:
	docker compose up -d

stop-infra:
	docker compose down

import:
	go run ./cmd/importcards/main.go
	