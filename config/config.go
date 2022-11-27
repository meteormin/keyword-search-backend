package config

import (
	"github.com/gofiber/fiber/v2"
	loggerMiddleware "github.com/gofiber/fiber/v2/middleware/logger"
	jwtWare "github.com/gofiber/jwt/v3"
	rsGen "github.com/miniyus/go-fiber/pkg/rs256"
	"log"
	"os"
	"path"
	"strconv"
)

func app() fiber.Config {
	return fiber.Config{
		AppName: os.Getenv("APP_NAME"),
	}
}

func logger() loggerMiddleware.Config {

	return loggerMiddleware.Config{
		Format:     "[${time}] ${ip}:${port} | (${pid}) ${status} - ${method} ${path}\n",
		TimeZone:   os.Getenv("TIME_ZONE"),
		TimeFormat: "2006-01-02 15:04:05",
	}
}

type DB struct {
	Host        string
	Dbname      string
	Username    string
	Password    string
	Port        string
	TimeZone    string
	SSLMode     bool
	AutoMigrate bool
}

func database() DB {
	autoMigrate, err := strconv.ParseBool(os.Getenv("DB_AUTO_MIGRATE"))

	if err != nil {
		autoMigrate = false
	}

	return DB{
		Host:        os.Getenv("DB_HOST"),
		Dbname:      os.Getenv("DB_DATABASE"),
		Username:    os.Getenv("DB_USERNAME"),
		Password:    os.Getenv("DB_PASSWORD"),
		Port:        os.Getenv("DB_PORT"),
		TimeZone:    os.Getenv("TIME_ZONE"),
		SSLMode:     false,
		AutoMigrate: autoMigrate,
	}
}

type Auth struct {
	Jwt jwtWare.Config
}

func auth() Auth {
	_, err := os.Stat(GetPath().DataPath)
	if err != nil {
		log.Fatalf("data path is not exists... %v", err)
	}

	secretPath := path.Join(GetPath().DataPath, "secret")

	_, err = os.Stat(secretPath)
	if err != nil {
		e := os.Mkdir(secretPath, os.FileMode(0755))
		if e != nil {
			log.Fatalf("%v", e)
		}
		log.Println("generate JWT secret keys...")
		rsGen.Generate(secretPath, 4096)
	}

	privateKey := path.Join(secretPath, "private.pem")

	priKey := rsGen.PrivatePemDecode(privateKey)

	return Auth{
		jwtWare.Config{
			SigningMethod: "RS256",
			SigningKey:    priKey.Public(),
			TokenLookup:   "header:Authorization",
		},
	}
}

func test() map[string]any {
	return map[string]any{
		"test_api": true,
		"name":     "TEST",
	}
}

type Configs struct {
	AppEnv   string
	AppPort  int
	Locale   string
	App      fiber.Config
	Logger   loggerMiddleware.Config
	Database DB
	Path     Path
	Auth     Auth
	Test     map[string]any
}

func GetConfigs() *Configs {
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))

	if err != nil {
		log.Printf("App Port is not numeric... %v", err)
		port = 8000
	}

	return &Configs{
		AppEnv:   os.Getenv("APP_ENV"),
		AppPort:  port,
		Locale:   os.Getenv("LOCALE"),
		App:      app(),
		Logger:   logger(),
		Database: database(),
		Path:     GetPath(),
		Auth:     auth(),
		Test:     test(),
	}
}

func InjectConfigContext(ctx *fiber.Ctx) error {
	ctx.Locals(Config, GetConfigs())

	return ctx.Next()
}
