package routes

import (
	"github.com/miniyus/keyword-search-backend/app"
	"github.com/miniyus/keyword-search-backend/auth"
	configure "github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/database"
	"github.com/miniyus/keyword-search-backend/internal/api_auth"
	"github.com/miniyus/keyword-search-backend/internal/group_detail"
	"github.com/miniyus/keyword-search-backend/internal/groups"
	"github.com/miniyus/keyword-search-backend/internal/host_search"
	"github.com/miniyus/keyword-search-backend/internal/hosts"
	"github.com/miniyus/keyword-search-backend/internal/jobs"
	"github.com/miniyus/keyword-search-backend/internal/search"
	"github.com/miniyus/keyword-search-backend/internal/short_url"
	"github.com/miniyus/keyword-search-backend/internal/test_api"
	"github.com/miniyus/keyword-search-backend/internal/users"
	"github.com/miniyus/keyword-search-backend/permission"
	"github.com/miniyus/keyword-search-backend/pkg/jwt"
	rsGen "github.com/miniyus/keyword-search-backend/pkg/rs256"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
	"github.com/miniyus/keyword-search-backend/utils"
	"gorm.io/gorm"
	"path"
)

const ApiPrefix = "/api"

func Api(apiRouter app.Router, a app.Application) {
	var cfg *configure.Configs
	a.Resolve(&cfg)

	if cfg == nil {
		configs := configure.GetConfigs()
		cfg = &configs
	}

	var db *gorm.DB
	a.Resolve(&db)

	if db == nil {
		db = database.GetDB()
	}

	var jDispatcher worker.Dispatcher
	a.Resolve(&jDispatcher)

	privateKey := rsGen.PrivatePemDecode(path.Join(cfg.Path.DataPath, "secret/private.pem"))
	tokenGenerator := jwt.NewGenerator(privateKey, privateKey.Public(), cfg.Auth.Exp)

	permissions := permission.NewPermissionsFromConfig(cfg.Permission)
	permissionCollection := permission.NewPermissionCollection(permissions...)

	authMiddlewareParam := auth.MiddlewaresParameter{
		Cfg: cfg.Auth.Jwt,
		DB:  db,
	}

	hasPermParam := permission.HasPermissionParameter{
		DB:           db,
		DefaultPerms: permissionCollection,
		FilterFunc:   group_detail.FilterFunc,
	}

	apiRouter.Route(
		test_api.Prefix,
		test_api.Register(jDispatcher),
	).Name("api.test_api")

	apiRouter.Route(
		api_auth.Prefix,
		api_auth.Register(
			api_auth.New(
				db,
				tokenGenerator,
			),
			authMiddlewareParam,
		),
	).Name("api.auth")

	// 해당 라인 이후로는 auth middleware가 공통으로 적용된다.
	apiRouter.Middleware(auth.Middlewares(authMiddlewareParam, permission.HasPermission(hasPermParam))...)
	// job 메타 데이터에 user_id 추가
	apiRouter.Middleware(jobs.AddJobMeta(jDispatcher, db))

	apiRouter.Route(
		jobs.Prefix,
		jobs.Register(
			jobs.New(
				utils.RedisClientMaker(cfg.RedisConfig),
				jDispatcher,
			),
		),
	)

	apiRouter.Route(
		groups.Prefix,
		groups.Register(groups.New(db)),
	).Name("api.groups")

	apiRouter.Route(
		users.Prefix,
		users.Register(users.New(db)),
	).Name("api.users")

	apiRouter.Route(
		hosts.Prefix,
		hosts.Register(hosts.New(db)),
	).Name("api.hosts")

	apiRouter.Route(
		search.Prefix,
		search.Register(search.New(db)),
	).Name("api.search")

	hostSearchHandler := host_search.New(db, jDispatcher)

	apiRouter.Route(
		host_search.Prefix,
		host_search.Register(hostSearchHandler, hasPermParam),
	).Name("api.hosts.search")

	apiRouter.Route(
		short_url.Prefix,
		short_url.Register(short_url.New(
			db,
			utils.RedisClientMaker(cfg.RedisConfig),
		)),
	).Name("api.short_url")

}
