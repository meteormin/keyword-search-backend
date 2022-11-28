package config

import (
	"os"
	"strconv"
)

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
