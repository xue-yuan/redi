package middleware

import (
	"redi/config"
	"redi/constants"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

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
