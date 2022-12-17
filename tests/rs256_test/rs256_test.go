package rs256_test

import (
	"github.com/miniyus/go-fiber/pkg/rs256"
	"log"
	"os"
	"testing"
)

func TestRs256(t *testing.T) {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	rs256.Generate(currentPath, 4096)
}
