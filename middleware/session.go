package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
)

var sessionStore *session.Store

func CreateSessionStore() {
	// Create a SQLite storage instance
	storage := sqlite3.New(sqlite3.Config{
		Database: "./fiber_session.db", // Path to your SQLite database file
		Table:    "sessions",           // Table name for storing session data
	})

	sessionStore = session.New(session.Config{
		Storage: storage, // Use the SQLite storage
	})
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
