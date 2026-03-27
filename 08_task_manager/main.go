package main

import (
	"fmt"
	"net/http"
	"os"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	IsDone      bool   `json:"is_done"`
}

var tasks []Task

func handleAddTask(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprint(w, "Название задачи не может быть пустым")
		return
	}
	task := Task{
		ID:          len(tasks) + 1,
		Description: name,
		IsDone:      false,
	}
	tasks = append(tasks, task)
	fmt.Fprintf(w, "Задача %d добавлена\n", task.ID)

}
func handleListTasks(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Список задач: %v\n", tasks)

}
func handleDoneTask(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprint(w, "Название задачи не может быть пустым")
		return
	}
	for i := 0; i < len(tasks); i++ {
		if tasks[i].Description == name {
			tasks[i].IsDone = true
			fmt.Fprintf(w, "Задача %d выполнена\n", tasks[i].ID)
			return
		}
	}
	fmt.Fprintf(w, "Задача %s не найдена\n", name)
}
func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprint(w, "Название задачи не может быть пустым")
		return
	}
	for i := 0; i < len(tasks); i++ {
		if tasks[i].Description == name {
			deletedID := tasks[i].ID
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "Задача %d удалена\n", deletedID)
			return
		}
	}
	fmt.Fprintf(w, "Задача %s не найдена\n", name)
}
func handleExit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "До свидания!")
	os.Exit(0)

}
func main() {
	http.HandleFunc("/add", handleAddTask)
	http.HandleFunc("/list", handleListTasks)
	http.HandleFunc("/done", handleDoneTask)
	http.HandleFunc("/delete", handleDeleteTask)
	http.HandleFunc("/exit", handleExit)
	fmt.Println("Привет! Это менеджер задач.")
	http.ListenAndServe(":8080", nil)
}
