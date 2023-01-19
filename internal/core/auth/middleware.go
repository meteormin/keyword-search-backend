package auth

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwtWare "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	configure "github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/core/container"
	"github.com/miniyus/keyword-search-backend/internal/core/context"
	"go.uber.org/zap"
	"log"
	"strconv"
	"time"
)

// 인증 관련 공통 미들웨어 작성

// User
// context에 저장될 유저 정보 구조체
type User struct {
	Id        uint   `json:"id"`
	GroupId   *uint  `json:"group_id"`
	Role      string `json:"role"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	ExpiresIn *int64 `json:"expires_in"`
}

// Middlewares
// 미들웨어 슬라이스 리턴
// 인증 관련된 미들웨어 함수의 집합으로 이 함수에 등록된 순서대로 실행 가능
func Middlewares(fn ...fiber.Handler) []fiber.Handler {
	// 순서 중요함
	mws := []fiber.Handler{
		JwtMiddleware,  // check exists jwt
		GetUserFromJWT, // get user information from jwt
		AccessLogMiddleware,
		CheckExpired, // check expired jwt
	}

	if len(fn) != 0 {
		mws = append(mws, fn...)
	}

	return mws
}

// AccessLogMiddleware
// log 찍힐 때 user 정보 추가
func AccessLogMiddleware(c *fiber.Ctx) error {
	var logger *zap.SugaredLogger
	logger, err := context.GetContext[*zap.SugaredLogger](c, context.Logger)

	if err != nil {
		return err
	}

	start := time.Now()
	err = c.Next()
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

// GetUserFromJWT
// get user information from jwt token
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

// JwtMiddleware
// jwt 유효성 체크 미들웨어
func JwtMiddleware(c *fiber.Ctx) error {
	config, err := context.GetContext[*configure.Configs](c, context.Config)
	if err != nil {
		return err
	}

	middleware := newJwtMiddleware(config.Auth.Jwt)

	return middleware(c)
}

func newJwtMiddleware(config jwtWare.Config) fiber.Handler {
	jwtConfig := config
	jwtConfig.ErrorHandler = jwtError
	return jwtWare.New(jwtConfig)
}

// jwtError
// jwt 생성과 해독(? decode...) 관련 에러 핸들링
func jwtError(c *fiber.Ctx, err error) error {
	var status int

	if err.Error() == "Missing or malformed JWT" {
		status = fiber.StatusBadRequest
		return fiber.NewError(status, err.Error())
	}

	return fiber.NewError(status, err.Error())
}

// CheckExpired
// jwt 만료 기간 체크 미들웨어
func CheckExpired(c *fiber.Ctx) error {
	wrapper, err := context.GetContext[container.Container](c, context.Container)

	if err != nil {
		return err
	}

	tokenRepository := NewRepository(wrapper.Database())

	user, err := context.GetContext[*User](c, context.AuthUser)
	if err != nil {
		return err
	}

	token, err := tokenRepository.FindByUserId(user.Id)
	if err != nil {
		statusCode := fiber.StatusUnauthorized
		return fiber.NewError(statusCode, "Can't Find User From Database")
	}

	if token.ExpiresAt.Unix() < time.Now().Unix() {
		statusCode := fiber.StatusUnauthorized
		return fiber.NewError(statusCode, "JWT is expired")
	}

	return c.Next()
}

func GetAuthUser(c *fiber.Ctx) (*User, error) {
	user, err := context.GetContext[*User](c, context.AuthUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}
