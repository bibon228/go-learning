package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Inventory - это наша собственная структура данных для хранения информации о товаре.
// Теги `json:"name"` и `json:"count"` указывают, какими ключами эти поля станут,
// когда мы превратим структуру в формат JSON для отправки клиенту.
// Без тегов поля назывались бы с большой буквы (Name, Count).
type Inventory struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// Ware (Склад) - это мапа (словарь), где ключом является строка (например, "Apple"),
// а значением - указатель (*) на структуру товара Inventory.
// Указатель (*) означает, что мы храним не копию товара, а прямую ссылку на товар в памяти.
// Это позволяет легко обновлять количество товара без его пересоздания.
var Ware = map[string]*Inventory{"Apple": {Name: "Яблоки", Count: 100}}

// handleStatus - функция-обработчик для пути /status
// Она показывает текущее состояние товара на складе.
func handleStatus(w http.ResponseWriter, r *http.Request) {
	// r.URL.Query().Get("name") берет параметр "name" из ссылки.
	// Например, из /status?name=Apple мы получим строку "Apple".
	name := r.URL.Query().Get("name")

	// Пытаемся достать товар по имени (name) из нашей мапы Ware.
	// ok будет true, если товар есть на складе, и false - если его нет.
	item, ok := Ware[name]
	if ok {
		// Товар найден! Указываем браузеру, что сейчас пришлем ему данные в формате JSON.
		w.Header().Set("Content-Type", "application/json")
		// Берем наш товар (item), превращаем его в JSON и сразу отправляем в ответ (w).
		json.NewEncoder(w).Encode(item)
	} else {
		// Товар не найден. Отправляем сообщение об ошибке с нужным HTTP-статусом (404 NotFound).
		http.Error(w, "Продукт не найден", http.StatusNotFound)
	}
}

// handleEditCount - функция-обработчик для пути /edit
// Она изменяет количество существующего товара на складе.
func handleEditCount(w http.ResponseWriter, r *http.Request) {
	// Достаем имя товара и новое количество из ссылки
	// (например: /edit?name=Apple&amount=50)
	name := r.URL.Query().Get("name")
	amountStr := r.URL.Query().Get("amount") // amountStr - это пока что текст, например "50"

	// Превращаем текст amountStr в реальное число (int).
	// Если пользователь ввел вместо числа буквы (amount=много), strconv.Atoi вернет ошибку (err)
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		// Зафиксировали ошибку. Сообщаем пользователю и прерываем выполнение функции (return).
		fmt.Fprint(w, "Колличество(amount) не может быть пустым или содержать буквы")
		return
	}

	// Ищем товар в мапе Ware
	item, ok := Ware[name]
	if ok {
		// Товар есть на складе. Так как у нас в мапе лежат указатели (*Inventory),
		// изменение item.Count обновит данные прямо на "складе" Ware для всех остальных запросов.
		item.Count = amount
		// Отправляем в ответ (пишем в браузер) текстовое сообщение об успешном изменении.
		fmt.Fprintf(w, "Колличество для %s было изменено на %d\n", name, amount)

		// Отправляем обновленный товар в формате JSON.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item)
	} else {
		// Товара с таким именем нет, изменить не получилось.
		fmt.Fprint(w, "На складе не найден товар с названием: %s", name)
	}
}

func main() {
	// Привязываем функции-обработчики к путям в адресной строке (маршрутизация).
	// Если зайти на /status, выполнится handleStatus.
	http.HandleFunc("/status", handleStatus)
	// Если зайти на /edit, выполнится handleEditCount.
	http.HandleFunc("/edit", handleEditCount)

	// Печатаем сообщение в терминал (это увидишь только ты, не пользователь сайта).
	// Очень удобно для отладки, чтобы понимать, что сервер успешно стартовал.
	fmt.Println("Сервер запущен! Попробуйте зайти на http://localhost:8080/status?name=Apple")

	// Слушаем порт 8080.
	// Эта функция "замораживает" программу и заставляет ее бесконечно ждать входящие запросы.
	http.ListenAndServe(":8080", nil)
}
