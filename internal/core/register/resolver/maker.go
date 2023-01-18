package resolver

import (
	goContext "context"
	"github.com/go-redis/redis/v9"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/core/container"
	cLogger "github.com/miniyus/keyword-search-backend/internal/core/logger"
	"github.com/miniyus/keyword-search-backend/internal/core/permission"
	"github.com/miniyus/keyword-search-backend/pkg/jwt"
	rsGen "github.com/miniyus/keyword-search-backend/pkg/rs256"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
	"go.uber.org/zap"
	"log"
	"path"
)

func MakeJwtGenerator(w container.Container) func() jwt.Generator {
	dataPath := w.Config().Path.DataPath

	privateKey := rsGen.PrivatePemDecode(path.Join(dataPath, "secret/private.pem"))

	return func() jwt.Generator {
		return &jwt.GeneratorStruct{
			PrivateKey: privateKey,
			PublicKey:  privateKey.Public(),
			Exp:        w.Config().Auth.Exp,
		}
	}
}

func MakeLogger(w container.Container) func() *zap.SugaredLogger {
	loggerConfig := w.Config().CustomLogger
	return func() *zap.SugaredLogger {
		return cLogger.New(parseLoggerConfig(loggerConfig))
	}
}

func parseLoggerConfig(loggerConfig config.LoggerConfig) cLogger.Config {
	return cLogger.Config{
		TimeFormat: loggerConfig.TimeFormat,
		FilePath:   loggerConfig.FilePath,
		Filename:   loggerConfig.Filename,
		MaxAge:     loggerConfig.MaxAge,
		MaxBackups: loggerConfig.MaxBackups,
		MaxSize:    loggerConfig.MaxSize,
		Compress:   loggerConfig.Compress,
		TimeKey:    loggerConfig.TimeKey,
		TimeZone:   loggerConfig.TimeZone,
		LogLevel:   loggerConfig.LogLevel,
	}
}

func MakePermissionCollection(w container.Container) func() permission.Collection {
	cfg := w.Config().Permission

	permCfg := permission.NewPermissionsFromConfig(parsePermissionConfig(cfg))

	return func() permission.Collection {
		return permission.NewPermissionCollection(permCfg...)
	}
}

func parsePermissionConfig(permissionConfig []config.PermissionConfig) []permission.Config {
	var permCfg []permission.Config
	for _, cfg := range permissionConfig {
		permCfg = append(permCfg, permission.Config{
			Name:      cfg.Name,
			GroupId:   cfg.GroupId,
			Methods:   parseMethodConstants(cfg.Methods),
			Resources: cfg.Resources,
		})
	}

	return permCfg
}

func parseMethodConstants(methods []config.PermissionMethod) []permission.Method {
	var authMethods []permission.Method
	for _, method := range methods {
		authMethods = append(authMethods, permission.Method(method))
	}

	return authMethods
}

func MakeJobDispatcher(c container.Container) func() worker.Dispatcher {
	opts := c.Config().QueueConfig

	opts.Redis = MakeRedisClient(c)

	return func() worker.Dispatcher {
		return worker.NewDispatcher(opts)
	}
}

func MakeRedisClient(c container.Container) func() *redis.Client {
	return func() *redis.Client {
		client := redis.NewClient(c.Config().RedisConfig)
		pong, err := client.Ping(goContext.Background()).Result()
		log.Print(pong, err)
		return client
	}
}
