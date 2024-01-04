package auth

import (
	"errors"
	"fmt"

	"github.com/Galish/loyalty-system/internal/app/entity"
	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID string
}

func GenerateToken(secretKey string, user *entity.User) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&JWTClaims{
			UserID: user.ID,
		},
	)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(secretKey string, tokenString string) (*JWTClaims, error) {
	var claims JWTClaims

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	return &claims, nil
}
