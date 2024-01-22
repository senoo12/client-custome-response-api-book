package main

import (
	bookApps "client-response-api-book/apps/book"
	userApps "client-response-api-book/apps/user"
	common "client-response-api-book/common"

	"github.com/gofiber/fiber/v2"
)

func main()  {
	app := fiber.New()
	common.ConnectDatabase()

	userApps.UserRoutes(app)
	common.JWTMiddleware(app)
	bookApps.BookRoutes(app)

	app.Listen(":3000")
}