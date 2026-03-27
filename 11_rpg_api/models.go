package main

// Структура Героя
type Hero struct {
	Name  string `json:"name"`
	Class string `json:"class"` // "Warrior", "Mage", "Tank"
	HP    int    `json:"hp"`
	MP    int    `json:"mp"`
	Level int    `json:"level"`
	Gold  int    `json:"gold"`
	Dmg   int    `json:"dmg"`
}

// Глобальная "база данных" (мапа)
var heroes = make(map[string]*Hero)
