package auth

import (
	"github.com/gofiber/fiber/v2"
	jwtWare "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	configure "github.com/miniyus/go-fiber/config"
	"github.com/miniyus/go-fiber/internal/core/api_error"
	"time"
)

// 공통 미들웨어 작성

type User struct {
	Id        uint
	GroupId   uint
	Username  string
	Email     string
	CreatedAt time.Time
}

func GetUserFromJWT(c *fiber.Ctx) error {
	jwtData, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.Next()
	}
	claims := jwtData.Claims.(jwt.MapClaims)

	userId := uint(claims["id"].(float64))
	groupId := uint(claims["group"].(float64))
	username := claims["username"].(string)
	createdAt := claims["createdAt"].(time.Time)

	currentUser := &User{
		Id:        userId,
		GroupId:   groupId,
		Username:  username,
		CreatedAt: createdAt,
	}

	c.Locals("authUser", currentUser)
	return c.Next()
}

func JwtMiddleware(c *fiber.Ctx) error {
	config, ok := c.Locals(configure.Config).(*configure.Configs)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Can not found Config Context...")
	}

	middleware := newJwtMiddleware(config.Auth.Jwt)

	return middleware(c)
}

func newJwtMiddleware(config jwtWare.Config) fiber.Handler {
	jwtConfig := config
	jwtConfig.ErrorHandler = jwtError
	return jwtWare.New(jwtConfig)
}

func jwtError(c *fiber.Ctx, err error) error {
	var errRes api_error.ErrorResponse

	if err.Error() == "Missing or malformed JWT" {
		errRes = api_error.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		}

		return errRes.Response(c)
	}

	errRes = api_error.ErrorResponse{
		Status:  "error",
		Code:    fiber.StatusUnauthorized,
		Message: "Invalid or expired JWT!",
	}

	return errRes.Response(c)
}
