package auth

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwtWare "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	configure "github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/core/api_error"
	"github.com/miniyus/keyword-search-backend/internal/core/context"
	"github.com/miniyus/keyword-search-backend/internal/core/database"
	"github.com/miniyus/keyword-search-backend/internal/entity"
	"go.uber.org/zap"
	"log"
	"strconv"
	"time"
)

// 공통 미들웨어 작성

type User struct {
	Id        uint   `json:"id"`
	GroupId   *uint  `json:"group_id"`
	Role      string `json:"role"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	ExpiresIn *int64 `json:"expires_in"`
}

func Middlewares() []fiber.Handler {
	// 순서 중요함
	mws := []fiber.Handler{
		JwtMiddleware,  // check exists jwt
		GetUserFromJWT, // get user information from jwt
		AccessLogMiddleware,
		CheckExpired, // check expired jwt
	}

	return mws
}

func HasPerm(action ...entity.Action) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := c.Locals(context.AuthUser).(*User)
		if currentUser.Role == entity.Admin.RoleToString() {
			return c.Next()
		}

		return fiber.ErrForbidden
	}
}

func AccessLogMiddleware(c *fiber.Ctx) error {
	logger, ok := c.Locals(context.Logger).(*zap.SugaredLogger)
	if !ok {
		log.Print("Failed Load logger context")
		return fiber.NewError(fiber.StatusInternalServerError, "Failed Load logger context")
	}

	start := time.Now()
	err := c.Next()
	elapsed := time.Since(start).Milliseconds()
	cu, ok := c.Locals(context.AuthUser).(*User)
	userID := ""
	if !ok {
		userID = "guest"
	} else {
		userID = strconv.Itoa(int(cu.Id))
	}
	req := c.Path()
	method := c.Method()

	errString := ""
	if err != nil {
		errString = fmt.Sprintf("| %s", err.Error())
	}

	logger.Info(fmt.Sprintf("user: %4s | IP: %15s | %6dms | %s | %s %s",
		userID, c.IP(), elapsed, method, req, errString))

	return err
}

func GetUserFromJWT(c *fiber.Ctx) error {

	jwtData, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		log.Print("access guest")
		return c.Next()
	}

	claims := jwtData.Claims.(jwt.MapClaims)

	userId := uint(claims["user_id"].(float64))

	var groupId uint
	if claims["group_id"] != nil {
		groupId = uint(claims["group_id"].(float64))
	}

	role := claims["role"].(string)
	username := claims["username"].(string)
	email := claims["email"].(string)
	createdAt := claims["created_at"].(string)

	var expiresIn int64
	if claims["expires_in"] != nil {
		expiresIn = int64(claims["expires_in"].(float64))
	}

	layout := "2006-01-02T15:04:05Z07:00"
	createdAtTime, err := time.Parse(layout, createdAt)
	if err != nil {
		return err
	}

	currentUser := &User{
		Id:        userId,
		GroupId:   &groupId,
		Role:      role,
		Username:  username,
		Email:     email,
		CreatedAt: createdAtTime.Format("2006-01-02 15:04:05"),
		ExpiresIn: &expiresIn,
	}

	c.Locals(context.AuthUser, currentUser)
	return c.Next()
}

func JwtMiddleware(c *fiber.Ctx) error {
	config, ok := c.Locals(context.Config).(*configure.Configs)
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
		errRes = api_error.NewErrorResponse(c, fiber.StatusBadRequest, err.Error())

		return errRes.Response()
	}

	errRes = api_error.NewErrorResponse(c, fiber.StatusBadRequest, "Invalid or expired JWT!")

	return errRes.Response()
}

func CheckExpired(c *fiber.Ctx) error {
	config, ok := c.Locals(context.Config).(*configure.Configs)
	if !ok {
		errRes := api_error.NewErrorResponse(c, fiber.StatusUnauthorized, "Can't Find Config Context")

		return errRes.Response()
	}

	tokenRepository := NewRepository(database.DB(config.Database))

	user, ok := c.Locals(context.AuthUser).(*User)
	if !ok {
		errRes := api_error.NewErrorResponse(c, fiber.StatusUnauthorized, "Can't Find User Context")

		return errRes.Response()
	}

	token, err := tokenRepository.FindByUserId(user.Id)
	if err != nil {
		errRes := api_error.NewErrorResponse(c, fiber.StatusUnauthorized, "Can't Find User From Database")

		return errRes.Response()
	}

	if token.ExpiresAt.Unix() > time.Now().Unix() {
		return c.Next()
	}

	errRes := api_error.NewErrorResponse(c, fiber.StatusUnauthorized, "JWT is expired")

	return errRes.Response()
}
