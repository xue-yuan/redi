package utils

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

func IsStructEmpty(s interface{}) bool {
	rv := reflect.ValueOf(s)
	if rv.Kind() != reflect.Struct {
		panic("isStructEmpty called with a non-struct type")
	}

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		if field.String() != "" {
			return false
		}
	}

	return true
}

func GetIP(c *fiber.Ctx) (ip string) {
	ip = c.Get("X-Forwarded-For")
	if ip == "" {
		return c.IP()
	}

	return ip
}
