.PHONY: build run test lint generate migrate db-up db-down dev

build:
	go build -o bin/app ./main.go

run:
	go run ./main.go

test:
	go test ./...

lint:
	golangci-lint run

generate:
	sqlc generate

migrate:
	goose -dir ./db/migrations postgres "$(DATABASE_URL)" up

db-up:
	docker compose up -d

db-down:
	docker compose down

dev:
	npm run dev
