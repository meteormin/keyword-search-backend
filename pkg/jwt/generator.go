package jwt

import (
	"crypto"
	"crypto/rsa"
	jwtLib "github.com/golang-jwt/jwt/v4"
	"log"
)

type Generator interface {
	Generate(claims jwtLib.Claims, privateKey *rsa.PrivateKey) (*string, error)
	GetPrivateKey() *rsa.PrivateKey
}

type GeneratorStruct struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  crypto.PublicKey
}

func (g *GeneratorStruct) Generate(claims jwtLib.Claims, privateKey *rsa.PrivateKey) (*string, error) {
	token := jwtLib.NewWithClaims(jwtLib.SigningMethodRS256, claims)
	t, err := token.SignedString(privateKey)
	if err != nil {
		log.Printf("token.SignedString: %v", err)
		return nil, err
	}

	return &t, nil
}

func (g *GeneratorStruct) GetPrivateKey() *rsa.PrivateKey {
	return g.PrivateKey
}
