package security

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtGenerator struct {
	jwtSecret   string
	jwtLifetime time.Duration
}

func NewJwtGenerator(jwtSecret string, jwtLifetimeHours int) *JwtGenerator {
	return &JwtGenerator{jwtSecret: jwtSecret, jwtLifetime: time.Hour * time.Duration(jwtLifetimeHours)}
}

func (g *JwtGenerator) GenerateJwtToken(login string) (string, error) {
	jwtWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   login,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.jwtLifetime)),
	})

	token, err := jwtWithClaims.SignedString([]byte(g.jwtSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}
