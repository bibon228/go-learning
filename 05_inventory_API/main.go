package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type Inventory struct {
	Name  string
	Count int
}

var Warehouse = map[string]*Inventory{
	"apple": {Name: "Apple", Count: 50},
}

func handleEditCount(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	amountStr := r.URL.Query().Get("amount")
	amount, _ := strconv.Atoi(amountStr)
	item, ok := Warehouse[name]
	if ok {
		item.Count += amount
		fmt.Fprintf(w, "Добавлено! Теперь %s на складе %d", name, item.Count)
	} else {
		fmt.Fprintf(w, "Товар %s не найден", name)
	}
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	item, ok := Warehouse[name]
	if ok {
		fmt.Fprintf(w, "Нашли! На складе %d %s", item.Count, name)
	} else {
		fmt.Fprint(w, "На складе нету такого товара!")
	}

}
func handleDelete(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprint(w, "Поле name не должно быть пустым")
		return
	}
	_, ok := Warehouse[name]
	if !ok {
		fmt.Fprintf(w, "Товар %s и так отсутствует на складе!", name)
		return
	}
	delete(Warehouse, name)
	fmt.Fprintf(w, "Товар %s успешно удален!", name)

}
func handleAdd(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("name")
	amountStr := r.URL.Query().Get("amount")
	amount, _ := strconv.Atoi(amountStr)
	if item == "" && amountStr == "" {
		fmt.Fprintf(w, "Имя не может быть пустым!")
		return

	}
	Warehouse[item] = &Inventory{
		Name:  item,
		Count: amount,
	}
	fmt.Fprintf(w, "Успех! Новый товар %s добавлен на склад в количестве %d", item, amount)
}

func main() {
	http.HandleFunc("/update", handleEditCount)
	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/delete", handleDelete)
	http.HandleFunc("/add", handleAdd)
	fmt.Println("Склад запущен на :8080")
	http.ListenAndServe(":8080", nil)
}
