package main

import (
	"log"
	"net/http"
	"personal-http-server/cmd/echo/custom_validator"
	"personal-http-server/cmd/echo/customer_error"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type User struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func main() {
	e := echo.New()
	e.Validator = &custom_validator.CustomValidator{Validator: validator.New()}

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello"})
	})

	e.POST("/user", func(c *echo.Context) error {
		u := new(User)

		if err := c.Bind(u); err != nil {
			return err
		}
		if err := c.Validate(u); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, u)
	})

	e.GET("/user:name", func(c *echo.Context) error {
		name := c.Param("name")
		return c.String(http.StatusOK, name)
	})

	e.GET("/cookie", func(c *echo.Context) error {
		writeCookie(c)
		return c.JSON(http.StatusOK, map[string]string{"message": "Cookie has been set"})
	})

	e.GET("/read-cookie", func(c *echo.Context) error {
		readCookie(c)
		return c.JSON(http.StatusOK, map[string]string{"message": "Cookie has been retrieved"})
	})

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)

	}
	e.HTTPErrorHandler = customer_error.CustomerErrorHandler
}

func writeCookie(c *echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = "g6i1o0"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	c.SetCookie(cookie)
}

func readCookie(c *echo.Context) error {
	for _, cookie := range c.Cookies() {
		log.Println(cookie.Name)
		log.Println(cookie.Value)
	}

	return nil

}
