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
	Name        string
	Driver      string
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

var connections map[string]*gorm.DB

func GetDB(name ...string) *gorm.DB {
	if len(name) == 0 {
		return connections["default"]
	}

	return connections[name[0]]
}

func switchDriver(driver string) func(dsn string) gorm.Dialector {
	switch driver {
	case "postgres":
		return postgres.Open
	case "pgsql":
		return postgres.Open
	default:
		return postgres.Open
	}
}

// New
// gorm.DB 객체 생성 함수
func New(config ...Config) *gorm.DB {
	var cfg Config
	if len(config) != 0 {
		cfg = config[0]
	} else {
		panic("새 데이터베이스 연결을 생성하려면 database.Config(이)가 필요합니다.")
	}

	var sslMode string = "disable"
	if cfg.SSLMode {
		sslMode = "enable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Host, cfg.Username, cfg.Password, cfg.Dbname, cfg.Port, sslMode, cfg.TimeZone,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		cfg.Logger,
	)

	driver := switchDriver(cfg.Driver)

	db, err := gorm.Open(driver(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalf("Failed: Connect DB %v", err)
	}

	log.Println("Success: Connect DB")

	if cfg.AutoMigrate {
		migrations.Migrate(db)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed: Connect sqlDB %v", err)
	}

	sqlDB.SetConnMaxLifetime(cfg.MaxLifeTime)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)

	connections[cfg.Name] = db

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
