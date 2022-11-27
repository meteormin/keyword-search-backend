package main

import (
	jwtLib "github.com/golang-jwt/jwt/v4"
	"github.com/miniyus/go-fiber/config"
	"github.com/miniyus/go-fiber/pkg/jwt"
	rsGen "github.com/miniyus/go-fiber/pkg/rs256"
	"log"
	"path"
)

func main() {
	c := jwtLib.MapClaims{
		"id": "smyoo",
	}

	jwtGenerator := func() *jwt.GeneratorStruct {
		dataPath := config.GetPath().DataPath
		privateKey := rsGen.PrivatePemDecode(path.Join(dataPath, "secret/private.pem"))

		return &jwt.GeneratorStruct{
			PrivateKey: privateKey,
			PublicKey:  privateKey.Public(),
		}
	}

	jwtGen := jwtGenerator()
	t, err := jwtGen.Generate(c, jwtGen.GetPrivateKey())
	if err != nil {
		log.Fatal(err)
	}

	log.Println(*t)
	log.Println(len(*t))
}
