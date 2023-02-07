PHONY: start
start:
	go run ./cmd/gofiber/main.go

.PHONY: build
build:
	swag init --parseDependency --parseInternal -g cmd/gofiber/main.go --output api
	go build -o build/gofiber ./cmd/gofiber/main.go

.PHONY: migrate
migrate:
	go run ./cmd/migrate/main.go

.PHONY: deploy
deploy:
	/usr/bin/env bash ./scripts/deploy.sh