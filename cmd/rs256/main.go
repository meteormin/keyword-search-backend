package main

import (
	"github.com/miniyus/go-fiber/pkg/rs256"
	"log"
	"os"
)

func main() {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	rs256.Generate(currentPath, 4096)
}
