package main

import (
	"github.com/miniyus/gofiber"
	"github.com/miniyus/gofiber/api_error"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/routes"
	_ "github.com/miniyus/keyword-search-backend/api/gofiber"
	"github.com/miniyus/keyword-search-backend/config"
	ksRoutes "github.com/miniyus/keyword-search-backend/routes"
)

// @title keyword-search-backend Swagger API Documentation
// @version 1.0.0
// @description keyword-search-backend API
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
	cfg := config.GetConfigs()
	appConfig := cfg.App
	appConfig.FiberConfig.ErrorHandler = api_error.OverrideDefaultErrorHandler(appConfig.Env)

	a := gofiber.New(cfg)

	// register routes
	a.Route(routes.ApiPrefix, func(router app.Router, app app.Application) {
		routes.Api(router, app)
		ksRoutes.Api(router, app)
	}, "api")

	// print status
	a.Status()
	// run application
	a.Run()

}
