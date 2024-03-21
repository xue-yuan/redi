package router

import (
	"redi/controllers"
	"redi/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")
	api.Get("/heartbeat", controllers.Heartbeat)

	v1 := api.Group("/v1")
	{
		user := v1.Group("/user")
		user.Post("/register", controllers.Register)
		user.Post("/login", controllers.Login)
		user.Post("/logout", middleware.HardAuth(), controllers.Logout)
		user.Post("/refresh_token", middleware.HardAuth(), controllers.RefreshToken)
		// user.Get("/", middleware.HardAuth(), controllers.GetUser)
		user.Get("/", middleware.SoftAuth(), controllers.GetUser)

		url := v1.Group("url")
		url.Post("/", middleware.SoftAuth(), controllers.CreateShortURL)
	}
}
