package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("umico_secret_key")

func GenerateToken(userID int) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом или гаечным )
	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil

}