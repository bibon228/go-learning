package main

// Здесь живут структуры (Hero, CreateHeroRequest, MarketRequest и т.д.)

type Hero struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Class string `json:"class"`
	Level int    `json:"level"`
	Gold  int    `json:"gold"`
	Dmg   int    `json:"dmg"`
	HP    int    `json:"hp"`
}

type CreateHeroRequest struct {
	Name  string `json:"name"`
	Class string `json:"class"`
}

type FightRequest struct {
	Hero1 string `json:"hero1"`
	Hero2 string `json:"hero2"`
}
