package controllers

import "github.com/gofiber/fiber/v2"

func Heartbeat(c *fiber.Ctx) error {
	return c.JSON((map[string]bool{"is_alive": true}))
}
