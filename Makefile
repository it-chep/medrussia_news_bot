include .env

PHONY: install
install:
	go install gotest.tools/gotestsum@latest
	go install github.com/joho/godotenv/cmd/dotenv@latest

test:
	gotestsum --format=testname

migrations-infra:
	goose -dir ./migrations postgres "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up


local: migrations-infra
	dotenv go run --lgflags