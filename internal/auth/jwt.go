package auth

import (
	"errors"
	"fmt"

	repo "github.com/Galish/loyalty-system/internal/repository"

	"github.com/golang-jwt/jwt/v4"
)

const secretKey = "yvdUuY)HSX}?&b':8/9N5"

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID string
}

func GenerateToken(user *repo.User) (string, error) {
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

func ParseToken(tokenString string) (*JWTClaims, error) {
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
