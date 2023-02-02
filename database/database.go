package database

import (
	"fmt"
	"github.com/miniyus/keyword-search-backend/database/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	gormLogger "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Config struct {
	Host        string
	Dbname      string
	Username    string
	Password    string
	Port        string
	TimeZone    string
	SSLMode     bool
	AutoMigrate bool
	Logger      gormLogger.Config
	MaxIdleConn int
	MaxOpenConn int
	MaxLifeTime time.Duration
}

var db *gorm.DB

// DB
// gorm.DB 객체 생성 함수
func DB(config Config) *gorm.DB {

	if db != nil {
		return db
	}

	var sslMode string = "disable"
	if config.SSLMode {
		sslMode = "enable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Host, config.Username, config.Password, config.Dbname, config.Port, sslMode, config.TimeZone,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		config.Logger,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalf("Failed: Connect DB %v", err)
	}

	log.Println("Success: Connect DB")

	if config.AutoMigrate {
		migrations.Migrate(db)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed: Connect sqlDB %v", err)
	}

	sqlDB.SetConnMaxLifetime(config.MaxLifeTime)
	sqlDB.SetMaxIdleConns(config.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.MaxOpenConn)

	return db
}

// HandleResult
// db 실행 결과 handle
func HandleResult(rs *gorm.DB) (*gorm.DB, error) {
	if rs.Error != nil {
		return nil, rs.Error
	}

	return rs, nil
}
