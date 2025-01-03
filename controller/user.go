package controller

import (
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/usysrc/nibbu/model"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash verifies a password against a hashed password
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// the register page
func Register(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{}, "layout")
}

func RegisterUser(c *fiber.Ctx) error {
	var registerData model.RegisterData
	if err := c.BodyParser(&registerData); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		slog.Error(err.Error())
		return err
	}

	hashedPassword, err := hashPassword(registerData.Password)
	if err != nil {
		return c.Render("registerform", fiber.Map{})
	}
	registerData.Password = hashedPassword

	err = model.RegisterUser(registerData)
	if err != nil {
		return c.Render("registerform", fiber.Map{
			"ErrorDescription": "Could not register user: username already in use.",
		})
	}

	user, err := model.GetUserByName(registerData.Username)
	if err != nil {
		slog.Error(err.Error())
		return c.Render("loginform", fiber.Map{
			"LoginFailed": true,
		})
	}

	sess, ok := c.Locals("session").(*session.Session)
	if !ok {
		slog.Error("'session' not found in locals.")
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	sess.Set("userID", user.ID)
	sess.Save()
	user.LoggedIn = true

	CreateSubdomain(user.Username)

	c.Response().Header.Add("hx-redirect", "/")
	return nil
}

// the login page
func Login(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*model.User)
	if !ok {
		user = &model.User{}
	}

	return c.Render("login", fiber.Map{
		"User": user,
	}, "layout")
}

// handle login of user
func LoginUser(c *fiber.Ctx) error {
	var loginData model.LoginData
	if err := c.BodyParser(&loginData); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		slog.Error(err.Error())
		return err
	}

	user, err := model.GetUserByName(loginData.Username)
	if err != nil {
		slog.Error(err.Error())
		return c.Render("loginform", fiber.Map{
			"LoginFailed": true,
		})
	}

	if user.Username != loginData.Username || !checkPasswordHash(loginData.Password, user.Password) {
		return c.Render("loginform", fiber.Map{
			"LoginFailed": true,
		})
	}

	sess, ok := c.Locals("session").(*session.Session)
	if !ok {
		slog.Error("'User' not found in locals.")
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	sess.Set("userID", user.ID)
	sess.Save()
	user.LoggedIn = true
	c.Response().Header.Add("hx-redirect", "/")
	return nil
}

func Logout(c *fiber.Ctx) error {
	if sess, ok := c.Locals("session").(*session.Session); ok {
		sess.Destroy()
	}
	return c.Redirect("/")
}

type Subdomain string

func GetAllSubdomains() ([]Subdomain, error) {
	users, err := model.GetAllUsers()
	if err != nil {
		return nil, err
	}
	domains := []Subdomain{}
	for _, user := range users {
		domains = append(domains, Subdomain(user.Username))
	}
	return domains, err
}
