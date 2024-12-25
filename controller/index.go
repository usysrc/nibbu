package controller

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/usysrc/nibbu/model"
)

// write the index
func Index(c *fiber.Ctx) error {
	items, err := model.GetAllPosts()
	if err != nil {
		c.Status(http.StatusInternalServerError)
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

	err = c.Render("index", fiber.Map{
		"Items": items,
		"User":  user,
	}, "layout")
	if err != nil {
		slog.Error(err.Error())
	}
	return err
}
