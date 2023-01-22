package resolver

import (
	goContext "context"
	"github.com/go-redis/redis/v9"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/logger"
	"github.com/miniyus/keyword-search-backend/internal/permission"
	"github.com/miniyus/keyword-search-backend/pkg/jwt"
	rsGen "github.com/miniyus/keyword-search-backend/pkg/rs256"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
	"go.uber.org/zap"
	"log"
	"path"
)

type JwtGeneratorConfig struct {
	DataPath string
	Exp      int
}

func MakeJwtGenerator(cfg JwtGeneratorConfig) func() jwt.Generator {
	dataPath := cfg.DataPath

	privateKey := rsGen.PrivatePemDecode(path.Join(dataPath, "secret/private.pem"))

	return func() jwt.Generator {
		return &jwt.GeneratorStruct{
			PrivateKey: privateKey,
			PublicKey:  privateKey.Public(),
			Exp:        cfg.Exp,
		}
	}
}

func MakeLogger(cfg config.LoggerConfig) func() *zap.SugaredLogger {
	loggerConfig := cfg
	return func() *zap.SugaredLogger {
		return logger.New(parseLoggerConfig(loggerConfig))
	}
}

func parseLoggerConfig(loggerConfig config.LoggerConfig) logger.Config {
	return logger.Config{
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

func MakePermissionCollection(cfg []config.PermissionConfig) func() permission.Collection {
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

type JobDispatcherConfig struct {
	WorkerCfg worker.DispatcherOption
	RedisCfg  *redis.Options
}

func MakeJobDispatcher(cfg JobDispatcherConfig) func() worker.Dispatcher {
	opts := cfg.WorkerCfg

	opts.Redis = MakeRedisClient(cfg.RedisCfg)

	return func() worker.Dispatcher {
		return worker.NewDispatcher(opts)
	}
}

func MakeRedisClient(cfg *redis.Options) func() *redis.Client {
	return func() *redis.Client {
		client := redis.NewClient(cfg)
		pong, err := client.Ping(goContext.Background()).Result()
		log.Print(pong, err)
		return client
	}
}
