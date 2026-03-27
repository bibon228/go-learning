package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	InitDB()
	defer DB.Close()

	http.HandleFunc("/create_user", handleCreateUser)
	http.HandleFunc("/list_users", handleListUsers)
	http.HandleFunc("/list_tasks", handleListTasks)
	http.HandleFunc("/create_task", handleCreateTask)
	http.HandleFunc("/update_status", handleUpdateStatus)
	http.HandleFunc("/delete_task", handleDeleteTask)
	fmt.Println("Сервер запущен на 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
