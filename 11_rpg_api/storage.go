package main

import (
	"encoding/json"
	"os"
)

// Здесь будут жить функции для работы с файлами
func saveHeroes() {
	data, err := json.MarshalIndent(heroes, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile("heroes.json", data, 0644)
}
func loadHeroes() {
	data, err := os.ReadFile("heroes.json")
	if err != nil {
		return
	}
	json.Unmarshal(data, &heroes)
}
