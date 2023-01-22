package create_admin

import (
	"errors"
	"github.com/miniyus/keyword-search-backend/internal/app"
	"github.com/miniyus/keyword-search-backend/internal/database"
	"github.com/miniyus/keyword-search-backend/internal/entity"
	"github.com/miniyus/keyword-search-backend/internal/permission"
	"github.com/miniyus/keyword-search-backend/internal/resolver"
	"github.com/miniyus/keyword-search-backend/internal/utils"
	"gorm.io/gorm"
	"log"
	"time"
)

func existsAdmin(db *gorm.DB) bool {
	admin := &entity.User{}
	rs := db.Where(entity.User{Role: string(entity.Admin)}).Find(admin)
	rs, err := database.HandleResult(rs)
	if rs.RowsAffected == 0 {
		return false
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	log.Println("Skip create admin: already exists admin account")
	return true
}

func CreateAdmin(app app.Application) {
	db := app.DB()
	if existsAdmin(db) {
		return
	}
	configs := app.Config()
	permCollectionFn := resolver.MakePermissionCollection(configs.Permission)

	caCfg := configs.CreateAdmin

	username := caCfg.Username
	email := caCfg.Email
	password := caCfg.Password

	if !caCfg.IsActive || username == "" || password == "" || email == "" {
		log.Println("Skip create admin: account info is empty")
		return
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Println(err)
		return
	}

	now := time.Now()

	user := &entity.User{
		Username:        username,
		Password:        hashedPassword,
		Email:           email,
		Role:            string(entity.Admin),
		EmailVerifiedAt: &now,
	}

	permissions := permCollectionFn()

	entPerms := make([]entity.Permission, 0)

	permissions.For(func(perm permission.Permission, i int) {
		entPerms = append(entPerms, permission.ToPermissionEntity(perm))
	})

	group := &entity.Group{
		Name:        "Admin",
		Permissions: entPerms,
		Users:       []entity.User{*user},
	}

	rs := db.Debug().Create(group)
	_, err = database.HandleResult(rs)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Success create admin")
	return
}
