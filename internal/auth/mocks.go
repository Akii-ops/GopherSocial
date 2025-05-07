package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TestAuthenticator struct {
}

const secret = "test"

var TestClaims = jwt.MapClaims{
	"sub": int64(42),
	"exp": time.Now().Add(time.Hour).Unix(),

	"iss": "test-iss",
	"aud": "test-aud",
}

func (a *TestAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TestClaims)

	tokenString, _ := token.SignedString([]byte(secret))

	return tokenString, nil
}

func (a *TestAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}
