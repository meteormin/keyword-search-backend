package register

import (
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/core/container"
	"github.com/miniyus/keyword-search-backend/internal/core/logger"
	"github.com/miniyus/keyword-search-backend/internal/core/permission"
	"github.com/miniyus/keyword-search-backend/pkg/jwt"
	rsGen "github.com/miniyus/keyword-search-backend/pkg/rs256"
	"path"
)

func makeJwtGenerator(w container.Container) func() jwt.Generator {
	return func() jwt.Generator {
		dataPath := w.Config().Path.DataPath

		privateKey := rsGen.PrivatePemDecode(path.Join(dataPath, "secret/private.pem"))

		return &jwt.GeneratorStruct{
			PrivateKey: privateKey,
			PublicKey:  privateKey.Public(),
			Exp:        w.Config().Auth.Exp,
		}
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
