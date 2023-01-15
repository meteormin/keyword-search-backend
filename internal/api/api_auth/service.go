package api_auth

import (
	"github.com/gofiber/fiber/v2"
	jwtLib "github.com/golang-jwt/jwt/v4"
	"github.com/miniyus/keyword-search-backend/internal/api/users"
	"github.com/miniyus/keyword-search-backend/internal/core/auth"
	"github.com/miniyus/keyword-search-backend/internal/core/logger"
	"github.com/miniyus/keyword-search-backend/internal/entity"
	"github.com/miniyus/keyword-search-backend/internal/utils"
	"github.com/miniyus/keyword-search-backend/pkg/jwt"
	"time"
)

type Service interface {
	SignIn(signIn *SignIn) (*entity.AccessToken, error)
	SignUp(signUp *SignUp) (*SignUpResponse, error)
	ResetPassword(pk uint, passwordStruct *ResetPasswordStruct) (*entity.User, error)
	RevokeToken(pk uint, token string) (bool, error)
	logger.HasLogger
}

type ServiceStruct struct {
	repo           auth.Repository
	userRepo       users.Repository
	tokenGenerator jwt.Generator
	logger.HasLoggerStruct
}

func NewService(repo auth.Repository, userRepo users.Repository, generator jwt.Generator) Service {
	return &ServiceStruct{
		repo:            repo,
		userRepo:        userRepo,
		tokenGenerator:  generator,
		HasLoggerStruct: logger.HasLoggerStruct{Logger: userRepo.GetLogger()},
	}
}

func hashPassword(password string) (string, error) {
	hashPass, err := utils.HashPassword(password)

	return hashPass, err
}

func hashCheck(hashPass string, password string) bool {
	return utils.HashCheck(hashPass, password)
}

func (s *ServiceStruct) generateToken(user *entity.User, exp int64) (*string, error) {
	claims := jwtLib.MapClaims{
		"user_id":    user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"group_id":   user.GroupId,
		"role":       user.Role,
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
		createdAt := time.Now()
		expTime := time.Duration(s.tokenGenerator.GetExp())

		expiresAt := createdAt.Add(time.Second * expTime)

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

	if err != nil {
		return nil, err
	}

	if user == nil {
		hashed, err := hashPassword(up.Password)
		if err != nil {
			return nil, err
		}

		user, err = s.userRepo.Create(entity.User{
			Username: up.Username,
			Email:    up.Email,
			Password: hashed,
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
				ExpiresAt: utils.JsonTime(token.ExpiresAt),
				ExpiresIn: token.ExpiresAt.Unix(),
			},
		}
		return res, nil
	} else {
		return nil, fiber.NewError(fiber.StatusConflict, "User already exists")
	}
}

func (s *ServiceStruct) ResetPassword(pk uint, passwordStruct *ResetPasswordStruct) (*entity.User, error) {
	if passwordStruct.Password != passwordStruct.PasswordConfirm {
		return nil, fiber.NewError(fiber.StatusBadRequest, "패스워드와 패스워드확인 필드가 일치하지않습니다.")
	}

	hashed, err := hashPassword(passwordStruct.Password)
	if err != nil {
		return nil, err
	}

	rsUser, err := s.userRepo.Update(pk, entity.User{
		Password: hashed,
	})

	if err != nil {
		return nil, err
	}

	return rsUser, nil
}

func (s *ServiceStruct) RevokeToken(pk uint, token string) (bool, error) {
	user, err := s.userRepo.Find(pk)
	if err != nil {
		return false, err
	}

	entToken, err := s.repo.FindByToken(token)
	if err != nil {
		return false, err
	}

	if entToken.UserId == user.ID {
		rs, err := s.repo.Delete(entToken.ID)
		if err != nil {
			return false, err
		}
		return rs, nil
	}

	return false, nil
}
