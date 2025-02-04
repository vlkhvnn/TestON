include .env

MIGRATIONS_PATH = ./cmd/migrations

.PHONY: migrate-create
migrate-create:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(name)

.PHONY: migrate-up
migrate-up:
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_ADDR)" up

.PHONY: migrate-down
migrate-down:
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_ADDR)" down $(steps)