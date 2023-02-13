package tests

import (
	"github.com/miniyus/gofiber/utils"
	"log"
	"testing"
)

func TestBase64Encode(t *testing.T) {
	log.Print(utils.Base64UrlEncode("132"))

}
