package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var sessionStore *session.Store

func CreateSessionStore() {
	sessionStore = session.New()
}

// Middleware to initialize session
func SessionMiddleware(c *fiber.Ctx) error {
	sess, err := sessionStore.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create session"})
	}
	c.Locals("session", sess)
	return c.Next()
}

// Middleware to protect routes
func AuthMiddleware(c *fiber.Ctx) error {
	sess := c.Locals("session").(*session.Session)
	userID := sess.Get("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	return c.Next()
}
