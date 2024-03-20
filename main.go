package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/itsig0/tallytome/handler"
)

func main() {

	app := fiber.New(fiber.Config{
		ServerHeader: "TallyTome",
	})

	// compression baby
	app.Use(compress.New())

	app.Static("/", "./public")

	app.Get("/", handler.Index)

	tracker := app.Group("/hp-mana-tracker")
	tracker.Get("/", handler.Tracker)
	tracker.Post("/update", handler.TrackerUpdate)
	tracker.Post("/damage", handler.TrackerDamage)
	tracker.Get("/check", handler.CheckStore)
	// app.Get("/*", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello, World!")
	// })

	app.Listen(":3000")
}
