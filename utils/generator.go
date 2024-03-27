package utils

import (
	"crypto/sha256"
	"fmt"
	"redi/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const SHORT_URL_LENGTH = 5

func GenerateUUID(tag string) (string, error) {
	u, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s", tag, u), nil
}

func GenerateShortURL(s, userID string, x int) string {
	return Base62Encode(sha256.Sum256([]byte(s + userID)))[:SHORT_URL_LENGTH+x]
}

func GenerateToken(uid string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": uid,
		"exp":     time.Now().Add(config.Config.TokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.Config.SecretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}
