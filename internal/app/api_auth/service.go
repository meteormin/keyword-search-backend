package api_auth

import (
	"github.com/gofiber/fiber/v2"
	jwtLib "github.com/golang-jwt/jwt/v4"
	"github.com/miniyus/go-fiber/internal/app/users"
	"github.com/miniyus/go-fiber/internal/entity"
	"github.com/miniyus/go-fiber/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Service interface {
	SignIn(signIn *SignIn) (*entity.AccessToken, error)
	SignUp(signUp *SignUp) (*SignUpResponse, error)
}

type ServiceStruct struct {
	repo           Repository
	userRepo       users.Repository
	tokenGenerator jwt.Generator
}

func NewService(repo Repository, userRepo users.Repository, generator jwt.Generator) *ServiceStruct {
	return &ServiceStruct{
		repo:           repo,
		userRepo:       userRepo,
		tokenGenerator: generator,
	}
}

func hashPassword(password string) string {
	hashPass, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(hashPass)
}

func hashCheck(hashPass string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
	return err == nil
}

func (s *ServiceStruct) generateToken(user *entity.User, exp int64) (*string, error) {
	claims := jwtLib.MapClaims{
		"userId":     user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"group_id":   user.GroupId,
		"expires_in": exp,
	}

	return s.tokenGenerator.Generate(claims, s.tokenGenerator.GetPrivateKey())
}

func (s *ServiceStruct) SignIn(in *SignIn) (*entity.AccessToken, error) {
	user, err := s.userRepo.FindByUsername(in.Username)

	if err != nil {
		return nil, err
	}

	if user != nil {
		if !hashCheck(user.Password, in.Password) {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "비밀번호가 틀렸습니다.")
		}

		expiresAt := time.Now().Add(time.Hour * 24)
		token, err := s.generateToken(user, expiresAt.Unix())
		if err != nil {
			return nil, err
		}

		if *token == "" {
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed Generate Token")
		}

		accessToken := entity.AccessToken{
			Token:     *token,
			UserId:    user.ID,
			ExpiresAt: expiresAt,
		}

		t, err := s.repo.Create(accessToken)
		if err != nil {
			return nil, err
		}

		return t, nil
	}

	return nil, fiber.NewError(fiber.StatusNotFound, "user not exists")
}

func (s *ServiceStruct) SignUp(up *SignUp) (*SignUpResponse, error) {
	user, err := s.userRepo.FindByEntity(entity.User{
		Username: up.Username,
		Email:    up.Email,
	})

	if user == nil {

		user, err = s.userRepo.Create(entity.User{
			Username: up.Username,
			Email:    up.Email,
			Password: hashPassword(up.Password),
		})
		if err != nil {
			return nil, err
		}

		if user == nil {
			return nil, fiber.NewError(fiber.StatusConflict, "Can not Create User...")
		}

		token, err := s.SignIn(&SignIn{
			Username: up.Username,
			Password: up.Password,
		})

		if err != nil {
			return nil, err
		}

		res := &SignUpResponse{
			UserId: user.ID,
			TokenInfo: TokenInfo{
				Token:     token.Token,
				ExpiresAt: token.ExpiresAt,
			},
		}
		return res, nil
	} else {
		return nil, fiber.NewError(fiber.StatusConflict, "User already exists")
	}
}
