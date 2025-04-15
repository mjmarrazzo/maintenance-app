package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mjmarrazzo/maintenance-app/internal/api"
	"github.com/mjmarrazzo/maintenance-app/internal/database"
	"github.com/mjmarrazzo/maintenance-app/internal/handlers"
)

func init() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Println("Warning: No .env file found")
	}
}

func main() {

	db, err := database.InitPool()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(api.ErrorMiddleware())

	categoryHandler := handlers.NewCategoryHandler(db)
	categoryHandler.RegisterRoutes(e)

	locationHandler := handlers.NewLocationHandler(db)
	locationHandler.RegisterRoutes(e)

	taskHandler := handlers.NewTaskHandler(db)
	taskHandler.RegisterRoutes(e)

	e.Static("/public", "public")

	e.Logger.Debug("Server starting on port 1323...")
	e.Logger.Fatal(e.Start(":1323"))
}
