package controller

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/usysrc/nibbu/model"
)

func ShowBlog(c *fiber.Ctx) error {
	posts, err := model.GetAllPosts()
	if err != nil {
		slog.Error(err.Error())
	}
	return c.Render("blog", fiber.Map{
		"Posts": posts,
	}, "layout")
}
