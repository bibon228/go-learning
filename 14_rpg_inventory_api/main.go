package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	InitDB()
	defer DB.Close()

	// Регистрация новых роутов
	http.HandleFunc("/buy_item", handleBuyItem)
	http.HandleFunc("/inventory", handleInventory)

	fmt.Println("Инвентарный сервер запущен на 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
