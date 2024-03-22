package handler

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/itsig0/tallytome/view/hptracker"
)

func Tracker(c fiber.Ctx) error {
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

func TrackerUpdate(c fiber.Ctx) error {

	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	data := hptracker.TrackerData{
		HP:                  c.FormValue("hp"),
		HPBase:              c.FormValue("hp"),
		HPStartPercentage:   "0",
		HPPercentage:        "100",
		Mana:                c.FormValue("mana"),
		ManaBase:            c.FormValue("mana"),
		ManaRegen:           c.FormValue("manaregen"),
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
	pre := "tracker_"
	sess.Set(pre+"HP", data.HP)
	sess.Set(pre+"HPBase", data.HPBase)
	sess.Set(pre+"HPPercentage", data.HPPercentage)
	sess.Set(pre+"HPStartPercentage", data.HPStartPercentage)

	sess.Set(pre+"Mana", data.Mana)
	sess.Set(pre+"ManaBase", data.ManaBase)
	sess.Set(pre+"ManaRegen", data.ManaRegen)
	sess.Set(pre+"ManaStartPercentage", data.ManaStartPercentage)
	sess.Set(pre+"ManaPercentage", data.ManaPercentage)

	// sess.Set("tracker", data)

	sess.Save()

	return render(c, hptracker.TrackerColumn(data))
}

func TrackerDamage(c fiber.Ctx) error {

	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	damage, err := strconv.Atoi(c.FormValue("damageInput"))
	if err != nil {
		c.SendStatus(418)
		return err
	} else if damage == 0 {
		c.SendStatus(418)
		return c.SendString("HP NOT NULL")
	}

	if damage < 0 {
		damage *= -1
	}

	currentHP, err := strconv.Atoi(fmt.Sprint(sess.Get("tracker_HP")))
	if err != nil {
		c.SendStatus(418)
		return err
	}

	baseHP, err := strconv.Atoi(fmt.Sprint(sess.Get("tracker_HPBase")))
	if err != nil {
		c.SendStatus(418)
		return err
	}

	heal := string(c.FormValue("heal"))
	savingThrow := string(c.FormValue("savingthrow"))

	if savingThrow == "on" || heal == "false" {
		damageFloat := float64(damage) * 0.33333
		damage -= int(math.Round(damageFloat))
	}

	if heal == "true" && damage > 0 {
		damage *= -1
	}

	newHP := currentHP - damage

	if newHP < 0 {
		newHP = 0
	}

	if newHP > baseHP {
		newHP = baseHP
	}

	newPercentage := (newHP * 100) / baseHP

	data := hptracker.TrackerData{
		HP:                fmt.Sprint(newHP),
		HPBase:            fmt.Sprint(sess.Get("tracker_HPBase")),
		HPStartPercentage: fmt.Sprint(sess.Get("tracker_HPPercentage")),
		HPPercentage:      fmt.Sprint(newPercentage),
	}

	sess.Set("tracker_HP", newHP)
	sess.Set("tracker_HPPercentage", newPercentage)

	sess.Save()

	return render(c, hptracker.HPTracker(data))
}

func CheckStore(c fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}
	log.Info(sess.Keys())
	return nil
}
