package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/itsig0/tallytome/view/login"
)

func Login(c fiber.Ctx) error {
	return render(c, login.Show())
}
