package config

import (
	jwtWare "github.com/gofiber/jwt/v3"
	mConfig "github.com/miniyus/gofiber/config"
	rsGen "github.com/miniyus/gofiber/pkg/rs256"
	"log"
	"os"
	"path"
)

func auth() mConfig.Auth {
	_, err := os.Stat(getPath().DataPath)
	if err != nil {
		log.Fatalf("data path is not exists... %v", err)
	}

	secretPath := path.Join(getPath().DataPath, "secret")

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

	return mConfig.Auth{
		Jwt: jwtWare.Config{
			SigningMethod: "RS256",
			SigningKey:    priKey.Public(),
			TokenLookup:   "header:Authorization,query:token",
		},
		Exp: 86400,
	}
}
