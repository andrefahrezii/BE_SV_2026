.PHONY: build run migrate-up migrate-down docker-up docker-down docker-logs test lint

APP_NAME=sv-backend
DOCKER_COMPOSE=docker-compose.yml

build:
	go build -o bin/$(APP_NAME) ./cmd/server

run:
	go run ./cmd/server

migrate-up:
	go run ./cmd/server -migrate up

migrate-down:
	echo "Drop schema manually with psql if needed:"

docker-up:
	docker compose -f $(DOCKER_COMPOSE) up --build -d

docker-down:
	docker compose -f $(DOCKER_COMPOSE) down

docker-logs:
	docker compose -f $(DOCKER_COMPOSE) logs -f app

test:
	go test ./...

lint:
	golangci-lint run
