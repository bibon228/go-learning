package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	InitDB()
	http.HandleFunc("/user", handleCreateUser)
	http.HandleFunc("/issue_card", handleIssueCard)
	http.HandleFunc("/top_up", handleTopUp)
	http.HandleFunc("/link_card", handleLinkCard)

	fmt.Println("Сервер запущен на 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
	// Входная точка в микросервис биллинга
}
