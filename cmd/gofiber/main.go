package main

import (
	"github.com/miniyus/go-fiber/internal/core"
)

// @title go-fiber Swagger API Documentation
// @version 0.0.1
// @description go-fiber API
// @contact.name miniyus
// @contact.url https://miniyus.github.io
// @contact.email miniyu97@gmail.com
// @license.name MIT
// @host localhost:9090
// @BasePath /
// @schemes http
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				   Bearer token type
func main() {
	core.Run()
}
