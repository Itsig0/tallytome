package handler

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

var store = session.New()

func render(c fiber.Ctx, component templ.Component) error {
	// or templ wil bork...
	c.Set("Content-type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}
