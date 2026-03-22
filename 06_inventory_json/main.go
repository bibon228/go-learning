package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Inventory struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

var Ware = map[string]*Inventory{"Apple": {Name: "Яблоки", Count: 100}}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	item, ok := Ware[name]
	if ok {

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item)
	} else {
		http.Error(w, "Продукт не найден", http.StatusNotFound)
	}
}
func handleEditCount(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	amountStr := r.URL.Query().Get("amount")
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		fmt.Fprint(w, "Колличество(amount) не может быть пустым")
		return
	}
	item, ok := Ware[name]
	if ok {
		item.Count = amount
		fmt.Fprintf(w, "Колличество для %s было изменено на %d\n", name, amount)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item)
	} else {
		fmt.Fprintf(w, "На складе не найден товар с названием: %s", name)

	}
}

func main() {
	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/edit", handleEditCount)
	fmt.Println("Сервер запущен!")
	http.ListenAndServe(":8080", nil)
}
