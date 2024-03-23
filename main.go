package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
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
	tracker.Post("/mana", handler.TrackerMana)
	tracker.Get("/newround", handler.TrackerRound)
	tracker.Get("/reset", handler.TrackerReset)
	// app.Get("/*", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello, World!")
	// })

	app.Listen(":3000")
}
