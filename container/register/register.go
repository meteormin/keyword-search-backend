package register

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/miniyus/go-fiber/config"
	"github.com/miniyus/go-fiber/container"
	"github.com/miniyus/go-fiber/internal/api_error"
	"github.com/miniyus/go-fiber/pkg/jwt"
	rsGen "github.com/miniyus/go-fiber/pkg/rs256"
	router "github.com/miniyus/go-fiber/routes"
	"path"
)

// Boot is High Priority
func boot(w container.Container) {
	w.App().Use(func(ctx *fiber.Ctx) error {
		ctx.Locals("configs", w.Config())
		return ctx.Next()
	})

	w.Inject("app", w.App())
	w.Inject("config", w.Config())
	w.Inject("db", w.Database())

	jwtGenerator := func() *jwt.GeneratorStruct {
		dataPath := config.GetPath().DataPath
		privateKey := rsGen.PrivatePemDecode(path.Join(dataPath, "secret/private.pem"))

		return &jwt.GeneratorStruct{
			PrivateKey: privateKey,
			PublicKey:  privateKey.Public(),
		}
	}
	var tg jwt.Generator
	w.Bind(&tg, jwtGenerator)
	w.Inject("jwtGenerator", jwtGenerator())
}

// Middlewares register middleware
func middlewares(w container.Container) {
	w.App().Use(logger.New(w.Config().Logger))
	w.App().Use(recover.New())
	w.App().Use(api_error.ErrorHandler)
}

// Routes register Routes
func routes(w container.Container) {
	router.SetRoutes(w)
}

func Resister(w container.Container) {
	boot(w)
	middlewares(w)
	routes(w)
}
