package resolver

import (
	"github.com/miniyus/keyword-search-backend/internal/core/container"
	"github.com/miniyus/keyword-search-backend/pkg/jwt"
	"go.uber.org/zap"
)

func TokenGenerator(c container.Container) jwt.Generator {
	var tokenGenerator jwt.Generator
	c.Resolve(&tokenGenerator)
	return tokenGenerator
}

func Logger(c container.Container) *zap.SugaredLogger {
	var logger *zap.SugaredLogger
	c.Resolve(&logger)
	return logger
}
