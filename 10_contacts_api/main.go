package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Структура Contact
type Contact struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// Глобальная мапа для хранения контактов
var contacts = map[string]*Contact{}

func handleAddContact(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	phone := r.URL.Query().Get("phone")
	if name == "" || phone == "" {
		fmt.Fprint(w, "Поля Name и Phone не могут быть пустыми")
		return
	}
	contacts[name] = &Contact{
		ID:    len(contacts) + 1,
		Name:  name,
		Phone: phone,
	}
	data, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		fmt.Fprint(w, "Ошибка при сохранении контактов")
		return
	}
	os.WriteFile("data.json", data, 0644)
	fmt.Fprintf(w, "Контакт %s добавлен\n", name)
}
func handleList(w http.ResponseWriter, r *http.Request) {
	list := []*Contact{}
	for _, c := range contacts {
		list = append(list, c)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		fmt.Fprint(w, "Ошибка при получении списка контактов")
		return
	}
	fmt.Fprint(w, string(data))
}
func handleDelContact(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprint(w, "Поле Name не может быть пустым")
		return
	}
	delete(contacts, name)
	data, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		fmt.Fprint(w, "Ошибка при удалении контакта")
		return
	}
	os.WriteFile("data.json", data, 0644)
	fmt.Fprintf(w, "Контакт %s удален\n", name)
}
func handleUpdateContact(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	phone := r.URL.Query().Get("phone")
	if name == "" || phone == "" {
		fmt.Fprint(w, "Поля Name и Phone не могут быть пустыми")
		return
	}
	contact, ok := contacts[name]
	if !ok {
		fmt.Fprint(w, "Контакт не найден")
		return
	}
	contact.Phone = phone
	data, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		fmt.Fprint(w, "Ошибка при обновлении контакта")
		return
	}
	os.WriteFile("data.json", data, 0644)
	fmt.Fprintf(w, "Контакт %s обновлен\n", name)
}
func main() {
	// 1. Попытаться прочитать контакты из файла при старте сервера
	data, err := os.ReadFile("data.json")
	if err == nil {
		json.Unmarshal(data, &contacts)
		fmt.Println("Контакты успешно загружены из файла!")
	} else {
		fmt.Println("Файл data.json не найден")
	}
	http.HandleFunc("/add", handleAddContact)
	http.HandleFunc("/list", handleList)
	http.HandleFunc("/delete", handleDelContact)
	http.HandleFunc("/update", handleUpdateContact)
	// 2. Инициализировать HTTP-обработчики

	fmt.Println("API Контактов запущено на порту 8080")
	http.ListenAndServe(":8080", nil)
}
