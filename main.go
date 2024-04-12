package main

import (
	"flag"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/itsig0/tallytome/handler"
)

func main() {
	var port int

	flag.IntVar(&port, "p", 3000, "Provide a port number")

	flag.Parse()

	app := fiber.New(fiber.Config{
		ServerHeader: "TallyTome",
	})

	// compression baby
	app.Use(compress.New())
	app.Use(logger.New())

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
	app.Get("/login", handler.Login)

	app.Listen(fmt.Sprint(":", port))
}
