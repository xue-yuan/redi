package middleware

import (
	"redi/config"
	"redi/constants"
	"redi/utils"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func parseToken(c *fiber.Ctx) (jwt.MapClaims, error) {
	auth := c.Get("Authorization")
	scheme, credentials := utils.GetAuthorizationSchemeAndParam(auth)
	if scheme == "" {
		return nil, constants.ErrJWTMissing
	}

	token, err := jwt.Parse(credentials, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, jwtware.ErrJWTMissingOrMalformed
	}

	return token.Claims.(jwt.MapClaims), nil
}

func Auth(t constants.AuthType) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := parseToken(c)
		if (err == constants.ErrJWTMissing) && (t == constants.SoftAuth) {
			return c.Next()
		} else if err != nil {
			return constants.BadRequestResponse(c, err)
		}

		exp, err := claims.GetExpirationTime()
		if err != nil {
			return constants.BadRequestResponse(c, err)
		}

		if exp.Time.Before(time.Now()) {
			return constants.BadRequestResponse(c, jwtware.ErrJWTMissingOrMalformed)
		}

		c.Locals(constants.UserID, claims["user_id"])
		return c.Next()
	}
}

func HardAuth() fiber.Handler {
	return Auth(constants.HardAuth)
}

func SoftAuth() fiber.Handler {
	return Auth(constants.SoftAuth)
}

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(config.Config.SecretKey)},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == jwtware.ErrJWTMissingOrMalformed.Error() {
		return constants.BadRequestResponse(c, err)
	}

	return constants.UnauthorizedResponse(c, err)
}
