package router

import (
	"redi/config"
	"redi/controllers"
	v1 "redi/controllers/v1"
	"redi/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Static("/image", config.Config.ImageFolder)
	app.Get("/heartbeat", controllers.Heartbeat)
	app.Get("/redirect/:short_url", controllers.RedirectURL)

	apiV1 := app.Group("/api/v1")
	{
		user := apiV1.Group("/user")
		{
			user.Post("/register", v1.Register)
			user.Post("/login", v1.Login)
			user.Post("/logout", middleware.HardAuth(), v1.Logout)
			user.Post("/refresh_token", middleware.HardAuth(), v1.RefreshToken)
			user.Get("/", middleware.HardAuth(), v1.GetUser)
		}

		url := apiV1.Group("/url")
		{
			url.Get("/list", middleware.HardAuth(), v1.GetShortURLs)
			url.Get("/", middleware.HardAuth(), v1.GetShortURL)
			url.Post("/", middleware.SoftAuth(), v1.CreateShortURL)
			url.Delete("/", middleware.HardAuth(), v1.DeleteShortURL)
			url.Post("/upload_image", middleware.HardAuth(), v1.UploadImage)
			url.Post("/open_graph", middleware.HardAuth(), v1.CreateOpenGraph)
			url.Put("/open_graph", middleware.HardAuth(), v1.UpdateOpenGraph)
			url.Delete("/open_graph", middleware.HardAuth(), v1.DeleteOpenGraph)
			url.Post("/customization", middleware.HardAuth(), v1.CreateCustomizedShortURL)
		}

		stat := apiV1.Group("/stat")
		{
			stat.Get("/count", middleware.HardAuth(), v1.StatCount)
			stat.Get("/list", middleware.HardAuth(), v1.GetStats)
		}
	}
}
