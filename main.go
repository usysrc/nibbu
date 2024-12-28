package main

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
	slogfiber "github.com/samber/slog-fiber"

	"github.com/usysrc/nibbu/controller"
	"github.com/usysrc/nibbu/filter"
	"github.com/usysrc/nibbu/middleware"
	"github.com/usysrc/nibbu/model"
)

type Host struct {
	fiber *fiber.App
}

func main() {
	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")
	// Add the markdown filter
	engine.AddFuncMap(map[string]any{
		"markdown": filter.MarkdownFilter,
		"date":     filter.Date,
	})
	engine.Reload(true)

	// Start fiber
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	model.Connect()
	defer model.Close()
	subdomains, err := controller.GetAllSubdomains()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	hosts := map[string]*Host{}
	for _, sub := range subdomains {
		subdomain := fiber.New(fiber.Config{
			EnablePrintRoutes: true,
			Views:             engine,
		})
		// Serve static files
		subdomain.Static("/", "./public")

		subdomain.Get("/", controller.ShowBlog)
		hosts[string(sub)+".localhost:3000"] = &Host{subdomain}
	}
	defaultApp := setupDefaultApp(engine)
	hosts["localhost:3000"] = &Host{defaultApp}

	for host := range hosts {
		log.Info(host)
	}

	// Add the host routing
	app.Use(func(c *fiber.Ctx) error {
		host := hosts[c.Hostname()]
		if host == nil {
			return c.Render("404", fiber.Map{}, "layout")
		} else {
			host.fiber.Handler()(c.Context())
			return nil
		}
	})

	// Start server
	if err := app.Listen(":3000"); err != nil {
		slog.Error(err.Error())
	}
}

func setupDefaultApp(engine *html.Engine) *fiber.App {
	defaultApp := fiber.New(fiber.Config{
		Views:             engine,
		EnablePrintRoutes: true,
	})

	// Ignore favicon requests
	defaultApp.Use(favicon.New())

	// Serve static files
	defaultApp.Static("/", "./public")

	// Add structured logging middleware
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	defaultApp.Use(slogfiber.New(logger))

	// Add the session middleware
	middleware.CreateSessionStore()
	defaultApp.Use(middleware.SessionMiddleware)
	defaultApp.Use(middleware.UserMiddleware)

	// Define routes
	defaultApp.Get("/", controller.Index)
	defaultApp.Get("/login", controller.Login)
	defaultApp.Get("/posts/edit/:url", middleware.AuthMiddleware, controller.EditPost)
	defaultApp.Delete("/posts/delete/:id", middleware.AuthMiddleware, controller.DeletePost)
	defaultApp.Post("/posts", controller.CreatePost)
	defaultApp.Put("/posts", controller.UpdatePost)
	defaultApp.Post("/loginuser", controller.LoginUser)
	defaultApp.Post("/logout", controller.Logout)
	defaultApp.Get("/logout", controller.Logout)
	defaultApp.Get("/register", controller.Register)
	defaultApp.Post("/registeruser", controller.RegisterUser)
	defaultApp.Get("/posts/new", middleware.AuthMiddleware, controller.NewPost)
	defaultApp.Get("/posts", controller.Posts)

	// Add the 404 handler
	// defaultApp.Use(func(c *fiber.Ctx) error {
	// 	return c.Render("404", fiber.Map{}, "layout")
	// })

	return defaultApp
}
