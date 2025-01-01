package middleware

import (
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
	"github.com/usysrc/nibbu/model"
)

var SessionStore *session.Store

func CreateSessionStore() {
	// Create a SQLite storage instance
	storage := sqlite3.New(sqlite3.Config{
		Database: "./fiber_session.db", // Path to your SQLite database file
		Table:    "sessions",           // Table name for storing session data
	})

	SessionStore = session.New(session.Config{
		Storage: storage, // Use the SQLite storage
	})
}

// Middleware to initialize session
func Session(c *fiber.Ctx) error {
	session, err := SessionStore.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	if session.Fresh() {
		session.Regenerate()
	}
	c.Locals("session", session)
	return c.Next()
}

// Middleware to get the user
func User(c *fiber.Ctx) error {
	sess, ok := c.Locals("session").(*session.Session)
	if !ok {
		slog.Error("No session struct found in locals.")
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	userID := sess.Get("userID")
	user := &model.User{}
	if userID != nil {
		id, err := strconv.Atoi(userID.(string))
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		user, err = model.GetUserByID(id)
		if err != nil { // user not found
			slog.Error(err.Error())
			if user == nil {
				user = &model.User{}
			}
			user.LoggedIn = false
		} else {
			user.LoggedIn = true
		}
	}
	c.Locals("user", user)
	return c.Next()
}

// Middleware to protect routes
func Auth(c *fiber.Ctx) error {
	sess, ok := c.Locals("session").(*session.Session)
	if !ok {
		slog.Error("No session struct found in locals.")
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	userID := sess.Get("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	return c.Next()
}

// Returns a fiber.Handler to be used as Middleware to enable CSRF
func CreateCSRF() fiber.Handler {
	return csrf.New(csrf.Config{
		KeyLookup:      "form:_csrf",
		ContextKey:     "csrf",
		CookieName:     "csrf_",
		CookieSameSite: "Lax",
		CookieSecure:   false, // !TODO: set to true in production
		Session:        SessionStore,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			slog.Error(err.Error())
			return err
		},
	})
}
