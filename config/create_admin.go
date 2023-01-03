package config

import (
	"os"
	"strconv"
)

type CreateAdminConfig struct {
	Username string
	Password string
	Email    string
	IsActive bool
}

func createAdminConfig() CreateAdminConfig {
	createAdmin := os.Getenv("CREATE_ADMIN")
	isActive, err := strconv.ParseBool(createAdmin)
	if err != nil {
		isActive = false
	}

	return CreateAdminConfig{
		IsActive: isActive,
		Username: os.Getenv("CREATE_ADMIN_USERNAME"),
		Password: os.Getenv("CREATE_ADMIN_PASSWORD"),
		Email:    os.Getenv("CREATE_ADMIN_EMAIL"),
	}
}
