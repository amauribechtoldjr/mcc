dev:
	air

mig-add:
	goose -s create $(name) sql -dir "./internal/adapters/postgresql/migrations"

mig-up:
	goose up

sql-gen:
	sqlc generate

start-infra:
	docker compose up -d

stop-infra:
	docker compose down

	