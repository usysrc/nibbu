package controller

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/usysrc/nibbu/model"
)

// write the index
func Index(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*model.User)
	if !ok {
		slog.Error("User was not set in locals.")
	}

	err := c.Render("index", fiber.Map{
		"User": user,
	}, "layout")

	if err != nil {
		slog.Error(err.Error())
	}
	return err
}
