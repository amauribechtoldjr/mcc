dev:
	air

migrate-add:
	goose -s create $(name) sql -dir "./internal/adapter/repository/postgres/migrations"

migrate-up:
	goose up

migrate-down:
	goose down

sql-gen:
	sqlc generate

start-infra:
	docker compose up -d

stop-infra:
	docker compose down

	