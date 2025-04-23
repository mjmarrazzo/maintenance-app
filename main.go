package main

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mjmarrazzo/maintenance-app/handlers"
	"github.com/mjmarrazzo/maintenance-app/internal/api"
	"github.com/mjmarrazzo/maintenance-app/internal/database"
)

var store *sessions.CookieStore

func init() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
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
	e.Use(session.Middleware(store))

	homeHandler := handlers.NewHomeHandler()
	homeHandler.RegisterRoutes(e)

	authHandler := handlers.NewAuthHandler(db)
	authHandler.RegisterRoutes(e)

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
