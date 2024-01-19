package book

import (
	"github.com/gofiber/fiber/v2"
)

func BookRoutes(app *fiber.App){
	app.Route("/books", func(api fiber.Router){
		api.Get("/", GetAllBooks)
		api.Get("/:id", GetBookById)
		api.Post("/", CreateBook)
		api.Put("/:id", UpdateBook)
		api.Delete("/:id", DeleteBook)
	})
}