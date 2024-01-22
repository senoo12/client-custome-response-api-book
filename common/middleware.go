package common

import (
	"github.com/gofiber/fiber/v2"
    jwtware "github.com/gofiber/contrib/jwt"
)

func JWTMiddleware(app *fiber.App){
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))
}