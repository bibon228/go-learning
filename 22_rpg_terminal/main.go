package main

import (
	"fmt"
	"terminal_rpg/models"
)

func NewGlad(name string) *models.Gladiator {
	return &models.Gladiator{
		Name:  name,
		HP:    100,
		MaxHP: 100,
		DMG:   10,
	}
}
func TakeDmg(g *models.Gladiator, dmg int) {
	g.HP -= dmg
	if g.HP < 0 {
		g.HP = 0
		panic("Гладиатор " + g.Name + " погиб!")
	}
}

func heal(g *models.Gladiator, heal int) {
	g.HP += heal
	if g.HP > g.MaxHP {
		g.HP = g.MaxHP
	}
	fmt.Println("Гладиатор", g.Name, "вылечился на", heal, "HP")
}

func main() {
	arena := make(map[string]*models.Gladiator)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("-------------")
			fmt.Println("💥 ИГРА ОКОНЧЕНА:", err)
			fmt.Println("Попробуйте начать заново.")
		}
	}()
	fmt.Println("Гладиаторские бои в терминале!")
	for {
		var command string
		fmt.Scan(&command)
		switch command {
		case "new":
			var name string
			fmt.Scan(&name)
			arena[name] = NewGlad(name)
		case "dmg":
			var name string
			var dmg int
			fmt.Scan(&name, &dmg)
			TakeDmg(arena[name], dmg)
		case "heal":
			var name string
			var healAmount int
			fmt.Scan(&name, &healAmount)
			heal(arena[name], healAmount)
		case "exit":
			return
		}
	}
}
