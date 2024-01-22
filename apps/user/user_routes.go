package user

import (
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App){
	app.Route("/users", func(api fiber.Router){
		api.Get("/", GetAllUsers)
		api.Get("/:id", GetUserByID)
		api.Post("/register", Register)
		api.Post("/login", Login)
		api.Put("/:id", UpdateUser)
		api.Delete("/:id", DeleteUser)
	})
}
