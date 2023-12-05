package auth

import (
	userRepo "loyalty-system/internal/repository/user"

	"github.com/golang-jwt/jwt/v4"
)

const secretKey = "yvdUuY)HSX}?&b':8/9N5"

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID string
}

func GenerateToken(user *userRepo.User) (string, error) {
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
