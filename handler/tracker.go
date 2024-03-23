package handler

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/itsig0/tallytome/view/hptracker"
)

func Tracker(c fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	hx := len(c.GetReqHeaders()["Hx-Request"]) > 0

	data := hptracker.TrackerData{
		Round: "0",
	}

	fields := []string{"HP", "HPBase", "HPStartPercentage", "HPPercentage", "Mana", "ManaBase", "ManaRegen", "ManaStartPercentage", "ManaPercentage", "Round"}
	for _, field := range fields {
		if val := sess.Get("tracker_" + field); val != nil {
			reflect.ValueOf(&data).Elem().FieldByName(field).SetString(fmt.Sprint(val))
		}
	}

	if sess.Get("tracker_Round") == nil {
		sess.Set("tracker_Round", "0")
		sess.Save()
	}

	// prevent animation reset
	data.HPStartPercentage = data.HPPercentage
	data.ManaStartPercentage = data.ManaPercentage

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

	bd := hptracker.BaseData{}

	mistake := false

	if len(data.HP) > 5 || len(data.Mana) > 5 {
		bd.Errors = "Einfach nur zu viel."
		mistake = true
	}

	intManaRegen, _ := strconv.Atoi(data.ManaRegen)
	if intManaRegen > 100 && !mistake {
		bd.Errors = "Mehr als 100% geht nicht"
		mistake = true
	}

	if mistake {
		c.Status(422)
		c.Append("HX-Retarget", "#baseStats")
		c.Append("HX-Reswap", "outerHTML")

		return render(c, hptracker.BaseStats(data, bd))
	}

	pre := "tracker_"
	sess.Set(pre+"HP", data.HP)
	sess.Set(pre+"HPBase", data.HPBase)
	sess.Set(pre+"HPPercentage", data.HPPercentage)
	sess.Set(pre+"HPStartPercentage", data.HPPercentage)

	sess.Set(pre+"Mana", data.Mana)
	sess.Set(pre+"ManaBase", data.ManaBase)
	sess.Set(pre+"ManaRegen", data.ManaRegen)
	sess.Set(pre+"ManaStartPercentage", data.ManaPercentage)
	sess.Set(pre+"ManaPercentage", data.ManaPercentage)

	sess.Set(pre+"Round", "0")

	sess.Save()

	c.Append("HX-Trigger", "BaseUpdated")

	return render(c, hptracker.TrackerColumn(data))
}

func TrackerDamage(c fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	damage, err := strconv.Atoi(c.FormValue("damageInput"))
	if err != nil {
		c.Status(422)
		c.Append("HX-Retarget", "#damageInputs")
		c.Append("HX-Reswap", "outerHTML")
		return render(c, hptracker.Hp(hptracker.DamageData{
			Errors: "Einfach nur zu viel.",
		}))
	}

	if damage <= 0 {
		c.Status(422)
		c.Append("HX-Retarget", "#damageInputs")
		c.Append("HX-Reswap", "outerHTML")
		return render(c, hptracker.Hp(hptracker.DamageData{
			Errors: "Darf nicht 0 oder kleiner sein.",
		}))
	}

	currentHP, err := strconv.Atoi(fmt.Sprint(sess.Get("tracker_HP")))
	if err != nil {
		c.Status(422)
		c.Append("HX-Retarget", "#damageInputs")
		c.Append("HX-Reswap", "outerHTML")
		return render(c, hptracker.Hp(hptracker.DamageData{
			Errors:      "Standard Werte nicht gesetzt.",
			Values:      c.FormValue("damageInput"),
			SavingThrow: c.FormValue("savingthrow"),
		}))
	}

	// no error handling here because it's already done for currentHP
	baseHP, _ := strconv.Atoi(fmt.Sprint(sess.Get("tracker_HPBase")))

	heal := c.FormValue("heal")
	savingThrow := c.FormValue("savingthrow")

	if savingThrow == "on" || heal == "false" {
		damageFloat := float64(damage) * 0.33333
		damage -= int(math.Round(damageFloat))
	}

	if heal == "true" {
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

	c.Append("HX-Trigger", "HPUpdated")

	return render(c, hptracker.HPTracker(data))
}

func TrackerMana(c fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	manaData := hptracker.ManaData{}

	mana, err := strconv.Atoi(c.FormValue("manaInput"))
	if err != nil {
		c.Status(422)
		c.Append("HX-Retarget", "#manaInputs")
		manaData.Errors = "Einfach nur zu viel."
		return render(c, hptracker.Mana(manaData))
	}

	if mana <= 0 {
		c.Status(422)
		c.Append("HX-Retarget", "#manaInputs")
		manaData.Errors = "Darf nicht 0 oder kleiner sein."
		return render(c, hptracker.Mana(manaData))
	}

	currentMana, err := strconv.Atoi(fmt.Sprint(sess.Get("tracker_Mana")))
	if err != nil {
		c.Status(422)
		c.Append("HX-Retarget", "#manaInputs")
		manaData.Errors = "Standard Werte nicht gesetzt."
		manaData.Values = c.FormValue("manaInput")
		return render(c, hptracker.Mana(manaData))
	}

	baseMana, _ := strconv.Atoi(fmt.Sprint(sess.Get("tracker_ManaBase")))

	add := c.FormValue("add")

	if add == "true" {
		mana *= -1
	}

	newMana := currentMana - mana

	if newMana < 0 {
		newMana = 0
	}

	if newMana > baseMana {
		newMana = baseMana
	}

	newPercentage := (newMana * 100) / baseMana

	data := hptracker.TrackerData{
		Mana:                fmt.Sprint(newMana),
		ManaBase:            fmt.Sprint(sess.Get("tracker_ManaBase")),
		ManaStartPercentage: fmt.Sprint(sess.Get("tracker_ManaPercentage")),
		ManaPercentage:      fmt.Sprint(newPercentage),
	}

	sess.Set("tracker_Mana", newMana)
	sess.Set("tracker_ManaPercentage", newPercentage)

	sess.Save()

	c.Append("HX-Trigger", "ManaUpdated")

	return render(c, hptracker.ManaTracker(data))
}

func TrackerRound(c fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	data := hptracker.TrackerData{}

	fields := []string{"HP", "HPBase", "HPStartPercentage", "HPPercentage", "Mana", "ManaBase", "ManaRegen", "ManaStartPercentage", "ManaPercentage", "Round"}
	for _, field := range fields {
		if val := sess.Get("tracker_" + field); val != nil {
			reflect.ValueOf(&data).Elem().FieldByName(field).SetString(fmt.Sprint(val))
		}
	}

	mana, _ := strconv.Atoi(data.Mana)
	manaBase, _ := strconv.Atoi(data.ManaBase)
	manaRegen, _ := strconv.Atoi(data.ManaRegen)
	round, _ := strconv.Atoi(data.Round)

	if manaBase == 0 || manaRegen == 0 {
		c.Status(418)
		return render(c, hptracker.TrackerHeader(data))
	}

	mana += (manaRegen * manaBase) / 100

	if mana > manaBase {
		mana = manaBase
	}

	// hp stays the same for now
	data.HPStartPercentage = data.HPPercentage

	data.Mana = fmt.Sprint(mana)
	data.Round = fmt.Sprint(round + 1)
	data.ManaStartPercentage = data.ManaPercentage
	data.ManaPercentage = fmt.Sprint((mana * 100) / manaBase)

	sess.Set("tracker_Mana", data.Mana)
	sess.Set("tracker_ManaPercentage", data.ManaPercentage)
	sess.Set("tracker_ManaStartPercentage", data.ManaStartPercentage)
	sess.Set("tracker_Round", data.Round)

	sess.Save()

	return render(c, hptracker.TrackerHeader(data))
}

func TrackerReset(c fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}
	sess.Reset()
	return render(c, hptracker.Show(true, hptracker.TrackerData{Round: "0"}))
}
