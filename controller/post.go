package controller

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/usysrc/nibbu/model"
)

// add an item to the db
func AddPost(c *fiber.Ctx) error {
	slog.Debug(string(c.Body()))
	var newItem model.Post
	if err := c.BodyParser(&newItem); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		slog.Error(err.Error())
		return err
	}
	err := model.NewPost(newItem)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}
	err = ListPosts(c)
	return err
}

// list the items
func ListPosts(c *fiber.Ctx) error {
	items, err := model.GetAllPosts()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}
	err = c.Render("list", fiber.Map{
		"Items": items,
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

// render the write page
func Write(c *fiber.Ctx) error {

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

	return c.Render("write", fiber.Map{
		"User": user,
	}, "layout")
}
