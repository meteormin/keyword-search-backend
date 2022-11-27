package auth

import (
	"github.com/gofiber/fiber/v2"
	jwtWare "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/miniyus/go-fiber/internal/api_error"
	"time"
)

func JwtMiddleware(config jwtWare.Config) fiber.Handler {
	config.ErrorHandler = jwtError
	return jwtWare.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	var errRes api_error.ErrorResponse

	if err.Error() == "Missing or malformed JWT" {
		errRes = api_error.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		}

		return errRes.Response(c)
	}

	errRes = api_error.ErrorResponse{
		Code:    fiber.StatusUnauthorized,
		Message: "Invalid or expired JWT!",
	}

	return errRes.Response(c)
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
