package main

import (
	"github/s6352410016/go-fiber-gorm-rest-api-auth-jwt-postgresql/config"
	"github/s6352410016/go-fiber-gorm-rest-api-auth-jwt-postgresql/database"
	"github/s6352410016/go-fiber-gorm-rest-api-auth-jwt-postgresql/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()
	database.ConnectDB()

	app := fiber.New()
	routes.SetUpRoutes(app)

	app.Listen(":8080")
}
