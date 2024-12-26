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
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/usysrc/nibbu/model"
)

func customURLEncode(input string) string {
	// Replace spaces with dashes
	input = strings.ReplaceAll(input, " ", "-")
	// Encode the string
	return url.QueryEscape(input)
}

// add an item to the db
func CreatePost(c *fiber.Ctx) error {
	slog.Debug(string(c.Body()))
	var newPost model.Post
	if err := c.BodyParser(&newPost); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		slog.Error(err.Error())
		return err
	}

	sess := c.Locals("session").(*session.Session)
	userID := sess.Get("userID")
	user := &model.User{}
	if userID != nil {
		id, err := strconv.Atoi(userID.(string))
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		user, err = model.GetUserByID(id)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		user.LoggedIn = true
	}

	newPost.Author = user.Username
	newPost.Date = time.Now().Format("2006-01-02 15:04:05")
	newPost.URL = template.URL(customURLEncode(newPost.Name))

	err := model.NewPost(newPost)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}
	err = ListPosts(c)
	return err
}

func UpdatePost(c *fiber.Ctx) error {
	var post model.Post
	if err := c.BodyParser(&post); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		slog.Error(err.Error())
		return err
	}

	sess := c.Locals("session").(*session.Session)
	userID := sess.Get("userID")
	user := &model.User{}
	if userID != nil {
		id, err := strconv.Atoi(userID.(string))
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		user, err = model.GetUserByID(id)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		user.LoggedIn = true
	}

	post.Author = user.Username
	// newPost.Date = time.Now().Format("2006-01-02 15:04:05")
	post.URL = template.URL(customURLEncode(post.Name))

	err := model.UpdatePost(post)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}
	err = ListPosts(c)
	return err
}

// list the items
func ListPosts(c *fiber.Ctx) error {
	posts, err := model.GetAllPosts()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}
	err = c.Render("list", fiber.Map{
		"Items": posts,
	})
	if err != nil {
		slog.Error(err.Error())
	}
	return err
}

// write single item
func Single(c *fiber.Ctx) error {
	type Param struct {
		ID int `json:"id"`
	}
	param := Param{}
	if err := c.ParamsParser(&param); err != nil {
		slog.Error(err.Error())
		return err
	}
	item, err := model.GetPost(param.ID)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	err = c.Render("single", fiber.Map{
		"Item": item,
	})
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

// render the new post page
func NewPost(c *fiber.Ctx) error {
	sess := c.Locals("session").(*session.Session)
	userID := sess.Get("userID")
	user := &model.User{}
	if userID != nil {
		id, err := strconv.Atoi(userID.(string))
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		user, err = model.GetUserByID(id)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		user.LoggedIn = true
	}

	return c.Render("posts-new", fiber.Map{
		"User": user,
	}, "layout")
}

// edit a post
func EditPost(c *fiber.Ctx) error {
	sess := c.Locals("session").(*session.Session)
	userID := sess.Get("userID")
	user := &model.User{}
	if userID != nil {
		id, err := strconv.Atoi(userID.(string))
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		user, err = model.GetUserByID(id)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		user.LoggedIn = true
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
	sess := c.Locals("session").(*session.Session)
	userID := sess.Get("userID")
	user := &model.User{}
	if userID != nil {
		id, err := strconv.Atoi(userID.(string))
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		user, err = model.GetUserByID(id)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		user.LoggedIn = true
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
