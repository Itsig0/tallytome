package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/itsig0/tallytome/view/home"
)

func Index(c *fiber.Ctx) error {
	hx := len(c.GetReqHeaders()["Hx-Request"]) > 0
	return render(c, home.Show(hx))
}
