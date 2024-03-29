package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/gofiber/jobqueue"
	"github.com/miniyus/gofiber/jobs"
	"github.com/miniyus/gofiber/log"
	"github.com/miniyus/gofiber/pkg/jwt"
	"github.com/miniyus/gofiber/utils"
	worker "github.com/miniyus/goworker"
	configure "github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/auth"
	"github.com/miniyus/keyword-search-backend/internal/galleries"
	"github.com/miniyus/keyword-search-backend/internal/groups"
	"github.com/miniyus/keyword-search-backend/internal/host_search"
	"github.com/miniyus/keyword-search-backend/internal/hosts"
	"github.com/miniyus/keyword-search-backend/internal/permission"
	"github.com/miniyus/keyword-search-backend/internal/photos"
	"github.com/miniyus/keyword-search-backend/internal/search"
	"github.com/miniyus/keyword-search-backend/internal/short_url"
	"github.com/miniyus/keyword-search-backend/internal/tasks"
	"github.com/miniyus/keyword-search-backend/internal/test_api"
	"github.com/miniyus/keyword-search-backend/internal/users"
	repo "github.com/miniyus/keyword-search-backend/repo"
	"go.uber.org/zap"
	"gorm.io/gorm"
	goLog "log"
)

const ApiPrefix = "/api"

func Api(apiRouter app.Router, a app.Application) {
	var cfg *configure.Configs
	err := a.Resolve(&cfg)
	if err != nil {
		goLog.Println(err)
	}

	if cfg == nil {
		configs := configure.GetConfigs()
		cfg = &configs
	}

	var db *gorm.DB
	err = a.Resolve(&db)
	if err != nil {
		goLog.Println(err)
	}

	if db == nil {
		db = database.GetDB()
	}

	var jDispatcher worker.Dispatcher
	jDispatcher = jobqueue.GetDispatcher()

	if jDispatcher == nil {
		err = a.Resolve(&jDispatcher)
		if err != nil {
			goLog.Println(err)
		}
	}

	var zLogger *zap.SugaredLogger
	err = a.Resolve(&zLogger)
	if err != nil {
		goLog.Println(err)
	}

	if zLogger == nil {
		zLogger = log.GetLogger()
	}

	privateKey := cfg.Auth.PrivateKey
	tokenGenerator := jwt.NewGenerator(privateKey, privateKey.Public(), cfg.Auth.Exp)

	authHandler := auth.New(db, repo.NewUserRepository(db), tokenGenerator)
	apiRouter.Route(
		auth.Prefix,
		auth.Register(authHandler, cfg.Auth.Jwt),
	).Name("api.auth")

	hasPermission := permission.HasPermission(permission.HasPermissionParameter{
		DB:           db,
		DefaultPerms: cfg.Permission,
		FilterFunc:   nil,
	})

	groupsHandler := groups.New(db)
	apiRouter.Route(
		groups.Prefix,
		groups.Register(groupsHandler),
		auth.JwtMiddleware(cfg.Auth.Jwt), auth.Middlewares(), hasPermission(),
	).Name("api.groups")

	usersHandler := users.New(db)
	apiRouter.Route(
		users.Prefix,
		users.Register(usersHandler),
		auth.JwtMiddleware(cfg.Auth.Jwt), auth.Middlewares(), hasPermission(),
	).Name("api.users")

	apiRouter.Route(
		hosts.Prefix,
		hosts.Register(hosts.New(db)),
		auth.JwtMiddleware(cfg.Auth.Jwt), auth.Middlewares(), hasPermission(),
	).Name("api.hosts")

	apiRouter.Route(
		search.Prefix,
		search.Register(search.New(db)),
		auth.JwtMiddleware(cfg.Auth.Jwt), auth.Middlewares(), hasPermission(),
	).Name("api.search")

	hostSearchHandler := host_search.New(db, jDispatcher)

	apiRouter.Route(
		host_search.Prefix,
		host_search.Register(hostSearchHandler),
		auth.JwtMiddleware(cfg.Auth.Jwt), auth.Middlewares(),
		hasPermission(),
	).Name("api.hosts.search")

	apiRouter.Route(
		short_url.Prefix,
		short_url.Register(short_url.New(
			db,
			utils.RedisClientMaker(cfg.RedisConfig),
		)),
		auth.JwtMiddleware(cfg.Auth.Jwt), auth.Middlewares(), hasPermission(),
	).Name("api.short_url")

	makeRedisClient := utils.RedisClientMaker(cfg.RedisConfig)

	getAuthUserId := func(ctx *fiber.Ctx) (uint, error) {
		user, err := auth.GetAuthUser(ctx)
		if err != nil {
			return 0, err
		}

		return user.Id, nil
	}

	jobsHandler := jobs.New(makeRedisClient, getAuthUserId, jDispatcher, jobqueue.GetRepository())
	apiRouter.Route(
		jobs.Prefix,
		jobs.Register(jobsHandler),
		auth.JwtMiddleware(cfg.Auth.Jwt), auth.Middlewares(), hasPermission(),
	)

	apiRouter.Route(
		tasks.Prefix,
		tasks.Register(),
	)

	apiRouter.Route(
		galleries.Prefix,
		galleries.New(db),
		auth.JwtMiddleware(cfg.Auth.Jwt), auth.Middlewares(), hasPermission(),
	)
	apiRouter.Route(
		photos.Prefix,
		photos.New(db),
		auth.JwtMiddleware(cfg.Auth.Jwt), auth.Middlewares(), hasPermission(),
	)

	apiRouter.Route(test_api.Prefix, test_api.Register(jDispatcher, utils.RedisClientMaker(cfg.RedisConfig)()))
}
