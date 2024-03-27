package utils

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	jwt.Claims
	UserID string
}

func GetAuthorizationSchemeAndParam(v string) (string, string) {
	vs := strings.Split(v, " ")
	if len(vs) != 2 {
		return v, ""
	}

	return vs[0], vs[1]
}

func IsValidPassword(h, p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p)) == nil
}

func HashPassword(p string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
