package main

import (
	"fmt"
	"os"
	"os/signal"
	"redi/config"

	"redi/database"

	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := config.Initialize(); err != nil {
		fmt.Println()
	}
	fmt.Println(config.Config)

	if err := database.Initialize(); err != nil {
		fmt.Println()
	}
	defer database.Pool.Close()

	app := fiber.New()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Gracefully Shutdown")
		app.Shutdown()
	}()

	app.Get("/heartbeat", func(c *fiber.Ctx) error {
		return c.JSON((map[string]bool{"is_alive": true}))
	})

	v1 := app.Group("/v1")
	{
		v1.Get("/foo", func(c *fiber.Ctx) error {
			return c.JSON(map[string]string{"foo": "bar"})
		})
	}

	if err := app.Listen(":5278"); err != nil {
		fmt.Println("panic")
	}

	fmt.Println("Running cleanup tasks...")
}
