package security

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"time"
)

type JwtService struct {
	jwtSecret   []byte
	jwtLifetime time.Duration
}

type CustomClaims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

var (
	ErrInvalidToken = errors.New("invalid token")
)

func NewJwtService(jwtSecret []byte, jwtLifetimeHours int) *JwtService {
	return &JwtService{jwtSecret: jwtSecret, jwtLifetime: time.Hour * time.Duration(jwtLifetimeHours)}
}

func (g *JwtService) GenerateJwtToken(userID int) (string, error) {
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// TODO: добавить в Subject login ?
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.jwtLifetime)),
		},
		UserID: userID,
	}

	jwtWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtWithClaims.SignedString([]byte(g.jwtSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (g *JwtService) ValidateJwtToken(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return g.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		zap.L().Error("invalid token")
		return nil, ErrInvalidToken
	}

	return claims, nil
}
