package core

import (
	"github.com/joho/godotenv"
	"log"
)

// init
// load dotenv
func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("failed dotenv load")
	}
}
