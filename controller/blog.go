package controller

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
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
