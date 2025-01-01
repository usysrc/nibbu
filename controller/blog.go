package controller

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/html/v2"
	"github.com/usysrc/nibbu/filter"
	"github.com/usysrc/nibbu/model"
)

func ShowBlog(c *fiber.Ctx) error {
	posts, err := model.GetAllPublishedPostsFromUser(c.Locals("blog").(string))
	if err != nil {
		slog.Error(err.Error())
	}
	return c.Render("blog", fiber.Map{
		"Posts": posts,
	}, "layout")
}

func SingleBlogPost(c *fiber.Ctx) error {
	post, err := model.GetPostByUrl(c.Params("url"))
	if err != nil {
		slog.Error(err.Error())
	}
	upvotes, err := model.GetUpvotesByPostID(post.ID)
	if err != nil {
		slog.Error(err.Error())
	}
	return c.Render("blog-single", fiber.Map{
		"Post":    post,
		"Upvotes": upvotes,
	}, "layout")
}

type Host struct {
	Fiber *fiber.App
}

var Hosts map[string]*Host

func CreateSubdomain(name string) {
	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")
	// Add the markdown filter
	engine.AddFuncMap(map[string]any{
		"markdown": filter.MarkdownFilter,
		"date":     filter.Date,
	})
	engine.Reload(true)

	subApp := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
		Views:             engine,
		PassLocalsToViews: true,
	})

	// Ignore Favicon
	subApp.Use(favicon.New())

	// Add the blog name to the locals
	subApp.Use(func(c *fiber.Ctx) error {
		c.Locals("blog", name)
		return c.Next()
	})

	// Serve static files
	subApp.Static("/", "./public")

	subApp.Get("/", ShowBlog)
	subApp.Get("/:url", SingleBlogPost)

	Hosts[name+".localhost:3000"] = &Host{subApp}
}
