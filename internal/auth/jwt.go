package auth

import (
	repo "loyalty-system/internal/repository"

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
