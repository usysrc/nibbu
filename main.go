package main

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
	slogfiber "github.com/samber/slog-fiber"

	"github.com/usysrc/usystack/controller"
	"github.com/usysrc/usystack/filter"
	"github.com/usysrc/usystack/middleware"
	"github.com/usysrc/usystack/model"
)

func main() {
	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")
	// Add the markdown filter
	engine.AddFuncMap(map[string]any{
		"markdown": filter.MarkdownFilter,
	})
	engine.Reload(true)

	// Start fiber
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Ignore favicon requests
	app.Use(favicon.New())

	// Serve static files
	app.Static("/", "./public")

	// Add structured logging middleware
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	app.Use(slogfiber.New(logger))

	// Add the session middleware
	middleware.CreateSessionStore()
	app.Use(middleware.SessionMiddleware)

	// Define routes
	model.Connect()
	defer model.Close()
	app.Get("/", controller.Index)
	app.Get("/login", controller.Login)
	app.Get("/item/:id", middleware.AuthMiddleware, controller.Single)
	app.Post("/add-item", controller.AddItem)
	app.Post("/loginuser", controller.LoginUser)
	app.Post("/logout", controller.Logout)
	app.Get("/register", controller.Register)
	app.Post("/registeruser", controller.RegisterUser)

	// Start server
	if err := app.Listen(":3000"); err != nil {
		slog.Error(err.Error())
	}
}
