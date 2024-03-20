package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/itsig0/tallytome/view/hptracker"
)

func Tracker(c *fiber.Ctx) error {
	hx := len(c.GetReqHeaders()["Hx-Request"]) > 0

	data := hptracker.TrackerData{
		HP:                  "0",
		HPBase:              "0",
		HPStartPercentage:   "0",
		HPPercentage:        "0",
		Mana:                "0",
		ManaBase:            "0",
		ManaRegen:           "0",
		ManaStartPercentage: "0",
		ManaPercentage:      "0",
	}

	return render(c, hptracker.Show(hx, data))
}

func TrackerUpdate(c *fiber.Ctx) error {

	_, err := store.Get(c)
	if err != nil {
		return err
	}

	// Parse the form data
	formdata := c.Context().PostArgs()

	// Create an instance of TrackerData
	data := hptracker.TrackerData{
		HP:                  string(formdata.Peek("hp")),
		HPBase:              string(formdata.Peek("hp")),
		HPStartPercentage:   "0",
		HPPercentage:        "100",
		Mana:                string(formdata.Peek("mana")),
		ManaBase:            string(formdata.Peek("mana")),
		ManaRegen:           string(formdata.Peek("manaregen")),
		ManaStartPercentage: "0",
		ManaPercentage:      "100",
	}

	// why no work?
	// values := reflect.ValueOf(data)
	// types := values.Type()
	// prefix := "tracker_"
	// for i := 0; i < values.NumField(); i++ {
	// 	sess.Set(prefix+types.Field(i).Name, values.Field(i))
	// 	log.Info("test")
	// }
	// pre := "tracker_"
	// sess.Set(pre+"HP", data.HP)
	// sess.Set(pre+"HPBase", data.HPBase)
	// sess.Set(pre+"HPStartPercentage", data.HPStartPercentage)
	// sess.Set(pre+"HPStartPercentage", data.HPStartPercentage)

	// sess.Set(pre+"Mana", data.Mana)
	// sess.Set(pre+"ManaBase", data.ManaBase)
	// sess.Set(pre+"ManaRegen", data.ManaRegen)
	// sess.Set(pre+"ManaStartPercentage", data.ManaStartPercentage)
	// sess.Set(pre+"ManaPercentage", data.ManaPercentage)

	// sess.Save()

	return render(c, hptracker.TrackerColumn(data))
}

func TrackerDamage(c *fiber.Ctx) error {

	// formdata := c.Context().PostArgs()

	data := hptracker.TrackerData{
		HPStartPercentage: "100",
		HPPercentage:      "66",
	}

	// log.Info(data)

	return render(c, hptracker.HPTracker(data))
}

func CheckStore(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}
	log.Info(sess.Keys())
	return nil
}
