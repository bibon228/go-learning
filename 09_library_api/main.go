package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Твоя структура Book пойдет сюда

type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	IsAvailable bool   `json:"is_available"`
	Count       int    `json:"count"`
}

var books = map[string]*Book{}

func handleAddBook(w http.ResponseWriter, r *http.Request) {

	title := r.URL.Query().Get("title")
	if title == "" {
		fmt.Fprint(w, "Название книги не может быть пустым")
		return
	}
	author := r.URL.Query().Get("author")
	if author == "" {
		fmt.Fprint(w, "Автор книги не может быть пустым")
		return
	}
	countStr := r.URL.Query().Get("count")
	if countStr == "" {
		fmt.Fprint(w, "Количество книг не может быть пустым")
		return
	}
	count, err := strconv.Atoi(countStr)
	if err != nil {
		fmt.Fprint(w, "Количество книг должно быть числом")
		return
	}

	books[title] = &Book{
		ID:          len(books) + 1,
		Title:       title,
		Author:      author,
		IsAvailable: true,
		Count:       count,
	}
	fmt.Fprintf(w, "Книга %s добавлена\n", title)
}
func handleTakeBook(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprint(w, "Название книги не может быть пустым")
		return
	}
	book, ok := books[name]
	if !ok {
		fmt.Fprint(w, "Книга не найдена на складе")
		return
	}
	if book.Count == 0 {
		fmt.Fprint(w, "Книга закончилась")
		return
	}
	book.Count--
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)

}

func handleList(w http.ResponseWriter, r *http.Request) {
	list := []Book{}
	for _, book := range books {
		list = append(list, *book)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)

}
func handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprint(w, "Поле Name не может быть пустым")
		return
	}
	delete(books, name)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)

}
func handleReturnBook(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprint(w, "Поле Name не может быть пустым")
		return
	}
	book, ok := books[name]
	if !ok {
		fmt.Fprintf(w, "Книга с таким именем не найдена")
		return

	}
	book.Count++

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)

}
func main() {
	http.HandleFunc("/add", handleAddBook)
	http.HandleFunc("/take", handleTakeBook)
	http.HandleFunc("/list", handleList)
	http.HandleFunc("/delete", handleDeleteBook)
	http.HandleFunc("/return", handleReturnBook)

	fmt.Println("API Библиотеки запущено на порту 8080")
	http.ListenAndServe(":8080", nil)
}
