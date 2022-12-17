package core

import (
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("failed dotenv load")
	}
}
