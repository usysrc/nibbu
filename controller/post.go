package controller

import (
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/usysrc/nibbu/model"
)

func customURLEncode(input string) string {
	// Replace spaces with dashes
	input = strings.ReplaceAll(input, " ", "-")
	// Encode the string
	return url.QueryEscape(input)
}

// add a post to the db
func CreatePost(c *fiber.Ctx) error {
	slog.Debug(string(c.Body()))
	var newPost model.Post
	if err := c.BodyParser(&newPost); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		slog.Error(err.Error())
		return err
	}

	user, ok := c.Locals("user").(*model.User)
	if !ok {
		slog.Error("'User' not found in locals.")
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	newPost.Author = user.Username
	newPost.Date = time.Now().Format("2006-01-02 15:04:05")
	newPost.URL = template.URL(customURLEncode(newPost.Name))
	newPost.Published = "draft"

	err := model.NewPost(newPost)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}

	c.Response().Header.Add("hx-redirect", "/posts/edit/"+string(newPost.URL))

	return nil
}

func UpdatePost(c *fiber.Ctx) error {
	var post model.Post
	if err := c.BodyParser(&post); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		slog.Error(err.Error())
		return err
	}

	user, ok := c.Locals("user").(*model.User)
	if !ok {
		slog.Error("'User' not found in locals.")
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	post.Author = user.Username
	// post.Date = time.Now().Format("2006-01-02 15:04:05")
	post.URL = template.URL(customURLEncode(post.Name))

	err := model.UpdatePost(post)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}
	c.SendString("Saved successfully.")
	return nil
}

// list the posts
func ListPosts(c *fiber.Ctx) error {
	posts, err := model.GetAllPosts()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}
	err = c.Render("list", fiber.Map{
		"Posts": posts,
	})
	if err != nil {
		slog.Error(err.Error())
	}
	return err
}

// write single post
func Single(c *fiber.Ctx) error {
	type Param struct {
		ID int `json:"id"`
	}
	param := Param{}
	if err := c.ParamsParser(&param); err != nil {
		slog.Error(err.Error())
		return err
	}
	post, err := model.GetPost(param.ID)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	err = c.Render("single", fiber.Map{
		"Post": post,
	})
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

// render the new post page
func NewPost(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*model.User)
	if !ok {
		slog.Error("'User' not found in locals.")
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return c.Render("posts-new", fiber.Map{
		"User": user,
	}, "layout")
}

// edit a post
func EditPost(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*model.User)
	if !ok {
		slog.Error("'User' not found in locals.")
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	post, err := model.GetPostByUrl(c.Params("url"))
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return c.Render("posts-edit", fiber.Map{
		"User": user,
		"Post": post,
	}, "layout")
}

// render the posts page
func Posts(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*model.User)
	if !ok {
		slog.Error("'User' not found in locals.")
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	posts, err := model.GetAllPostsFromUser(user.Username)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}
	return c.Render("posts", fiber.Map{
		"User":  user,
		"Posts": posts,
	}, "layout")
}

func DeletePost(c *fiber.Ctx) error {
	idString := c.Params("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return err
	}
	slog.Info("id is " + strconv.Itoa(id))
	err = model.DeletePost(id)
	if err != nil {
		slog.Error(err.Error())
		return c.SendString("Something went wrong while deleting.")
	}
	c.Response().Header.Add("hx-redirect", "/posts")
	return nil
}

func PreviewPost(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*model.User)
	if !ok {
		slog.Error("'User' not found in locals.")
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	post, err := model.GetPostByUrl(c.Params("url"))
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return c.Render("preview", fiber.Map{
		"User": user,
		"Post": post,
	}, "layout")
}

func PublishPost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	err = model.PublishPost(id)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	c.Response().Header.Add("hx-refresh", "true")
	return c.SendString("Published!")
}

func UnpublishPost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	err = model.UnpublishPost(id)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	c.Response().Header.Add("hx-refresh", "true")
	return c.SendString("Published!")
}

func UpvotePost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	err = model.UpvotePost(id, c.IP())
	if err != nil {
		slog.Error(err.Error())
	}
	c.Response().Header.Add("hx-refresh", "true")
	return c.SendString("Upvoted!")
}
