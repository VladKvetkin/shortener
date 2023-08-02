package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JwtTokenBuilder struct{}

type claims struct {
	jwt.RegisteredClaims
	UserID string
}

const secretKey = "testsecretkey"
const tokenExp = time.Hour * 3

func (jwtb *JwtTokenBuilder) BuildJWTToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: uuid.NewString(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (jwtb *JwtTokenBuilder) GetUserID(tokenString string) (string, error) {
	claims := &claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		return "", err
	}

	if !token.Valid || claims.UserID == "" {
		return "", fmt.Errorf("token is not valid")
	}

	return claims.UserID, nil
}
