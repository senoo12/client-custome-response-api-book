package main

import (
	bookApps "client-response-api-book/apps/book"
	infra "client-response-api-book/infra"

	"github.com/gofiber/fiber/v2"
	// "github.com/go-sql-driver/mysql"
)

func main()  {
	app := fiber.New()
	infra.ConnectDatabase()

	bookApps.BookRoutes(app)

	app.Listen(":3000")
}