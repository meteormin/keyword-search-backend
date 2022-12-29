PHONY: start
start:
	go run ./cmd/gofiber/main.go

.PHONY: build
build:
	swag init --parseDependency --parseInternal --parseDepth 1 -g ./cmd/gofiber/main.go --output ./api/gofiber
	go build -o build/gofiber ./cmd/gofiber/main.go

.PHONY: migrate
migrate:
	go run ./cmd/migrate/main.go