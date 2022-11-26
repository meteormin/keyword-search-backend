.PHONY: start
start:
	go run ./cmd/gofiber/main.go

.PHONY: build
build:
	go build -o build/gofiber ./cmd/gofiber/main.go

.PHONY: migrate
migrate:
	go run ./cmd/migrate/main.go