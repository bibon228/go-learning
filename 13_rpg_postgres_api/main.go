package main

// Здесь запускается сервер и регистрируются роуты

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	InitDB()
	defer DB.Close()

	// Регистрация роутов
	http.HandleFunc("/create", handleCreate)
	http.HandleFunc("/list", handleListHeroes)
	http.HandleFunc("/levelup", handleLevelUp)
	http.HandleFunc("/market", handleMarket)
	http.HandleFunc("/fight", handleFight)

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
