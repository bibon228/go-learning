package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	InitDB()
	http.HandleFunc("/add", handleAddUser)
	http.HandleFunc("/create", handleCreateCard)
	http.HandleFunc("/topup", handleToPup)
	http.HandleFunc("/link", handleLinkCard)

	fmt.Println("Сервер запущен на 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
