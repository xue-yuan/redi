package main

import (
	"fmt"
	"os"
	"os/signal"
	"redi/config"
	"redi/database"
	"redi/middleware"
	"redi/redis"
	"redi/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := config.Initialize(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := database.Initialize(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer database.Pool.Close()

	if err := redis.Initialize(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer redis.Client.Close()

	app := fiber.New()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Gracefully Shutdown")
		app.Shutdown()
	}()

	app.Use(middleware.SetupContext)
	router.SetupRoutes(app)

	if err := app.Listen(":5278"); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Running cleanup tasks...")
}
