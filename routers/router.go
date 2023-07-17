package routers

import (
	"github.com/abe27/api/crypto/controllers"
	"github.com/gofiber/fiber/v2"
)

func Routers(c *fiber.App) {
	c.Get("/", controllers.HelloController)
	r := c.Group("/api/v1")
	r.Get("", controllers.HelloController)

	//// Router User
	user := r.Group("/member")
	user.Get("/me", controllers.ProfileController)
	user.Post("/login", controllers.LoginController)
	user.Post("/register", controllers.RegisterController)
	user.Delete("/logout", controllers.LogOutController)
}
