package main

import (
	"fmt"
	"net/http"
)

type Visitors = map[string]bool

var myVisitors = Visitors{
	"Sex":  true,
	"Flex": true,
}

func checkAccess(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	_, ok := myVisitors[name]
	if ok {
		fmt.Fprintf(w, "Доступ пользователю %s запрещен!\n", name)
	} else {
		fmt.Fprintf(w, "Добро пожаловать %s\n", name)
	}
}
func addToBlacklist(w http.ResponseWriter, r *http.Request) {
	newName := r.URL.Query().Get("newName")
	if newName == "" {
		fmt.Fprintf(w, "Введите имя!")
	} else {
		myVisitors[newName] = true
		fmt.Fprintf(w, "Пользователь %s добавлен в blacklist!", newName)
	}
}
func main() {
	http.HandleFunc("/add", addToBlacklist)
	http.HandleFunc("/check", checkAccess)
	fmt.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", nil)
}
