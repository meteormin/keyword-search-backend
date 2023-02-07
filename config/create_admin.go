package config

import (
	mConfig "github.com/miniyus/gofiber/config"
	"os"
	"strconv"
)

func createAdminConfig() mConfig.CreateAdminConfig {
	createAdmin := os.Getenv("CREATE_ADMIN")
	isActive, err := strconv.ParseBool(createAdmin)
	if err != nil {
		isActive = false
	}

	return mConfig.CreateAdminConfig{
		IsActive: isActive,
		Username: os.Getenv("CREATE_ADMIN_USERNAME"),
		Password: os.Getenv("CREATE_ADMIN_PASSWORD"),
		Email:    os.Getenv("CREATE_ADMIN_EMAIL"),
	}
}
