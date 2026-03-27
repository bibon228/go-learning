package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 1. Загружаем героев из файла
	loadHeroes()

	// 2. Регистрируем пути
	http.HandleFunc("/heroes", handleListHeroes)
	http.HandleFunc("/create", handleCreateHero)
	http.HandleFunc("/level_up", handleLevelUp)
	http.HandleFunc("/transfer_gold", handleTransferGold)
	http.HandleFunc("/fight", handleFight)
	http.HandleFunc("/market", handleMarket)

	// 3. Запускаем сервер
	fmt.Println("RPG Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", nil)
}
