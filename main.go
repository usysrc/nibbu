package main

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
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
		Views:             engine,
		PassLocalsToViews: true,
	})
	model.Connect()
	defer model.Close()
	subdomains, err := controller.GetAllSubdomains()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	hosts := map[string]*Host{}
	for _, name := range subdomains {
		subApp := fiber.New(fiber.Config{
			EnablePrintRoutes: true,
			Views:             engine,
			PassLocalsToViews: true,
		})

		// Ignore Favicon
		subApp.Use(favicon.New())

		// Add the blog name to the locals
		subApp.Use(func(c *fiber.Ctx) error {
			c.Locals("blog", string(name))
			return c.Next()
		})

		// Serve static files
		subApp.Static("/", "./public")

		subApp.Get("/", controller.ShowBlog)
		subApp.Get("/:url", controller.SingleBlogPost)

		hosts[string(name)+".localhost:3000"] = &Host{subApp}
	}
	defaultApp := setupDefaultApp(engine)
	hosts["localhost:3000"] = &Host{defaultApp}

	for host := range hosts {
		slog.Debug(host)
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
	if err := app.Listen("localhost:3000"); err != nil {
		slog.Error(err.Error())
	}
}

func setupDefaultApp(engine *html.Engine) *fiber.App {
	defaultApp := fiber.New(fiber.Config{
		Views:             engine,
		EnablePrintRoutes: true,
		PassLocalsToViews: true,
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
	defaultApp.Use(middleware.Session)
	defaultApp.Use(middleware.User)

	// Add the CSRF middleware
	csrfMiddleware := middleware.CreateCSRF()
	defaultApp.Use(csrfMiddleware)

	// Define routes
	defaultApp.Get("/", controller.Index)
	defaultApp.Get("/login", controller.Login)
	defaultApp.Get("/posts/edit/:url", middleware.Auth, controller.EditPost)
	defaultApp.Get("/posts/preview/:url", middleware.Auth, controller.PreviewPost)
	defaultApp.Delete("/posts/delete/:id", middleware.Auth, controller.DeletePost)
	defaultApp.Post("/posts", middleware.Auth, controller.CreatePost)
	defaultApp.Put("/posts", middleware.Auth, controller.UpdatePost)
	defaultApp.Post("/loginuser", controller.LoginUser)
	defaultApp.Post("/logout", middleware.Auth, controller.Logout)
	defaultApp.Get("/logout", middleware.Auth, controller.Logout)
	defaultApp.Get("/register", controller.Register)
	defaultApp.Post("/registeruser", controller.RegisterUser)
	defaultApp.Get("/posts/new", middleware.Auth, controller.NewPost)
	defaultApp.Get("/posts", middleware.Auth, controller.Posts)

	// Add the 404 handler
	// defaultApp.Use(func(c *fiber.Ctx) error {
	// 	return c.Render("404", fiber.Map{}, "layout")
	// })

	return defaultApp
}
