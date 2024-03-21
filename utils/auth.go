package utils

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
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
