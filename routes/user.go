package routes

import (
	"github/s6352410016/go-fiber-gorm-rest-api-auth-jwt-postgresql/handlers"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App) {
	user := app.Group("/api")
	user.Post("/signup", handlers.SignUp)
	user.Post("/signin", handlers.SignIn)

	user.Use("/profile", jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("AT_SECRET"))},
	}))
	user.Get("/profile", handlers.ShowProfile)

	user.Use("/refresh", jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("RT_SECRET"))},
	}))
	user.Post("/refresh", handlers.Refresh)
}
