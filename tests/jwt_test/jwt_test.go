package jwt_test

import (
	jwtLib "github.com/golang-jwt/jwt/v4"
	"github.com/miniyus/go-fiber/config"
	"github.com/miniyus/go-fiber/pkg/jwt"
	rsGen "github.com/miniyus/go-fiber/pkg/rs256"
	"path"
	"testing"
)

func TestJwt(t *testing.T) {
	c := jwtLib.MapClaims{
		"id": "smyoo",
	}

	jwtGenerator := func() *jwt.GeneratorStruct {
		dataPath := config.GetConfigs().Path.DataPath
		privateKey := rsGen.PrivatePemDecode(path.Join(dataPath, "secret/private.pem"))

		return &jwt.GeneratorStruct{
			PrivateKey: privateKey,
			PublicKey:  privateKey.Public(),
		}
	}

	jwtGen := jwtGenerator()
	token, err := jwtGen.Generate(c, jwtGen.GetPrivateKey())
	if err != nil {
		t.Errorf("Failed Generaate Token... %v", err)
	}

	t.Log(token)
}
