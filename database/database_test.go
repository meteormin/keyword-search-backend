package database_test

import (
	"github.com/miniyus/keyword-search-backend/database"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log"
	"testing"
	"time"
)

var defaultDB *gorm.DB

func TestNew(t *testing.T) {
	cfg := database.Config{
		Name:        "default",
		Driver:      "postgres",
		Host:        "localhost",
		Dbname:      "go_fiber",
		Username:    "smyoo",
		Password:    "smyoo",
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
		MaxOpenConn: 10,
		MaxLifeTime: time.Hour,
	}

	defaultDB = database.New(cfg)

	sql, err := defaultDB.DB()
	if err != nil {
		t.Error(err)
	}
	log.Print(sql.Stats())
}

func TestGetDB(t *testing.T) {
	get := database.GetDB()
	if get != defaultDB {
		t.Error(get)
	}

	sql, err := defaultDB.DB()
	if err != nil {
		t.Error(err)
	}
	log.Print(sql.Stats())

	sql2, err := get.DB()
	if err != nil {
		t.Error(err)
	}
	log.Println(sql2.Stats())
}

type Cnt struct {
	Cnt int
}

func TestHandleResult(t *testing.T) {
	get := database.GetDB()
	var cnt Cnt
	tx := get.Raw("SELECT count(*) as cnt FROM pg_catalog.pg_tables").Scan(&cnt)
	tx, err := database.HandleResult(tx)
	if err != nil {
		t.Error(err)
	}

	log.Print(cnt.Cnt)
	if cnt.Cnt != 77 {
		t.Error("not matched real db")
	}
}
