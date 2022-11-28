package register

import (
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/miniyus/go-fiber/config"
	"github.com/miniyus/go-fiber/internal/core/api_error"
	"github.com/miniyus/go-fiber/internal/core/auth"
	"github.com/miniyus/go-fiber/internal/core/container"
	"github.com/miniyus/go-fiber/internal/core/logger"
	router "github.com/miniyus/go-fiber/internal/routes"
	"github.com/miniyus/go-fiber/pkg/jwt"
	rsGen "github.com/miniyus/go-fiber/pkg/rs256"
	"go.uber.org/zap"
	"path"
)

// Boot is High Priority
func boot(w container.Container) {
	w.Inject("app", w.App())
	w.Inject("config", w.Config())
	w.Inject("db", w.Database())

	jwtGenerator := func() *jwt.GeneratorStruct {
		dataPath := w.Config().Path.DataPath
		privateKey := rsGen.PrivatePemDecode(path.Join(dataPath, "secret/private.pem"))

		return &jwt.GeneratorStruct{
			PrivateKey: privateKey,
			PublicKey:  privateKey.Public(),
		}
	}

	var tg jwt.Generator
	w.Bind(&tg, jwtGenerator)
	w.Resolve(&tg)
	w.Inject("jwtGenerator", &tg)

	var log *zap.SugaredLogger
	w.Bind(&log, logger.GetLogger)
	w.Resolve(&log)
	w.Inject("logger", &log)
}

// Middlewares register middleware
func middlewares(w container.Container) {
	w.App().Use(flogger.New(w.Config().Logger))
	w.App().Use(recover.New())
	w.App().Use(config.InjectConfigContext)
	w.App().Use(logger.Middleware)
	w.App().Use(api_error.ErrorHandler)
	w.App().Use(auth.GetUserFromJWT)
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
