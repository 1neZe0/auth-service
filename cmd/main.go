package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"os"

	"auth-service/config"
	"auth-service/controller"
	"auth-service/db"
)

func main() {
	config.LoadEnv()

	if err := db.ConnectDB(os.Getenv("DATABASE_URL")); err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	app := fiber.New()

	controller.SetupRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
