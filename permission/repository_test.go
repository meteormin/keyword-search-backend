package permission_test

import (
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/database"
	"github.com/miniyus/keyword-search-backend/permission"
	gormLogger "gorm.io/gorm/logger"
	"testing"
	"time"
)

func TestRepositoryStruct(t *testing.T) {
	var groupId uint = 1
	db := database.DB(config.DB{
		Host:        "localhost",
		Dbname:      "go_fiber",
		Username:    "",
		Password:    "",
		Port:        "5432",
		TimeZone:    "Asia/Seoul",
		SSLMode:     false,
		AutoMigrate: false,
		Logger: gormLogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormLogger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
		MaxIdleConn: 10,
		MaxOpenConn: 100,
		MaxLifeTime: time.Hour,
	})

	repo := permission.NewRepository(db)

	get, err := repo.Get(groupId)
	if err != nil {
		t.Error(err)
	}
	save, err := repo.Save(get)
	if err != nil {
		t.Error(err)
	}

	t.Error(save)
}
