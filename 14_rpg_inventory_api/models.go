package main

type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type BuyRequest struct {
	HeroName string `json:"hero_name"`
	ItemName string `json:"item_name"`
}
type Hero struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Class string `json:"class"`
	Level int    `json:"level"`
	Gold  int    `json:"gold"`
	Dmg   int    `json:"dmg"`
	HP    int    `json:"hp"`
}
