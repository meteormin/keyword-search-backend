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
	GetExp() int
}

type GeneratorStruct struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  crypto.PublicKey
	Exp        int
}

func NewGenerator(priv *rsa.PrivateKey, pub crypto.PublicKey, exp int) Generator {
	return &GeneratorStruct{
		priv, pub, exp,
	}
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

func (g *GeneratorStruct) GetExp() int {
	return g.Exp
}
