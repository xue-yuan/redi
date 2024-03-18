package router

import (
	"redi/controllers"
	"redi/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/heartbeat", controllers.Heartbeat)

	api := app.Group("/api")

	v1 := api.Group("/v1")
	{
		user := v1.Group("/user")
		user.Post("/register", controllers.Register)
		user.Post("/login", controllers.Login)
		user.Post("/logout", middleware.Protected(), controllers.Logout)
		user.Post("/refresh_token", middleware.Protected(), controllers.RefreshToken)
		user.Get("/", middleware.Protected(), controllers.GetUser)

		url := v1.Group("url")
		url.Post("", controllers.CreateShortURL)
	}
}
