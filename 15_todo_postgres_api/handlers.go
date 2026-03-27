package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var input CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Fatal("Не удалось прочитать пользователя из БД", err)
		return
	}
	_, err = DB.Exec("INSERT INTO users (name) VALUES ($1)", input.Name)
	if err != nil {
		log.Fatal("Не удалось добавить пользователя", err)
		return
	}
	fmt.Fprintf(w, "Пользователь %s успешно создан в БД!", input.Name)
}

func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var input CreateTaskRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Fatal("Не удалось прочитать задачу из БД", err)
		return
	}
	_, err = DB.Exec("INSERT INTO tasks (user_id, title) VALUES ($1, $2)", input.UserID, input.Title)
	if err != nil {
		log.Fatal("Не удалось добавить задачу", err)
		return
	}
	fmt.Fprintf(w, "Задача %s успешно создана в БД!", input.Title)
}

func handleListTasks(w http.ResponseWriter, r *http.Request) {
	UserID := r.URL.Query().Get("user_id")
	if UserID == "" {
		fmt.Fprint(w, "Укажите user_id")
		return
	}
	rows, err := DB.Query("SELECT id, title, is_done, user_id FROM tasks WHERE user_id = $1", UserID)
	if err != nil {
		log.Fatal("Не удалось получить список задач", err)
		return
	}
	defer rows.Close()
	list := []Task{}
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Done, &task.UserID)
		if err != nil {
			log.Fatal("Не удалось прочитать задачу", err)
			return
		}
		list = append(list, task)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func handleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	var input UpdateTaskStatusRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Fatal("Не удалось прочитать задачу из БД", err)
		return
	}
	_, err = DB.Exec("UPDATE tasks SET is_done = $1 WHERE id = $2", input.Done, input.TaskID)
	if err != nil {
		log.Fatal("Не удалось обновить задачу", err)
		return
	}
	fmt.Fprintf(w, "Задача %d успешно обновлена в БД!", input.TaskID)
}

func handleListUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal("Не удалось получить список пользователей", err)
		return
	}
	defer rows.Close()
	list := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			log.Fatal("Не удалось прочитать пользователя", err)
			return
		}
		list = append(list, user)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)

}

func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	TaskID := r.URL.Query().Get("task_id")
	if TaskID == "" {
		fmt.Fprint(w, "Укажите task_id")
		return
	}
	_, err := DB.Exec("DELETE FROM tasks WHERE id = $1", TaskID)
	if err != nil {
		log.Fatal("Не удалось удалить задачу", err)
		return
	}
	fmt.Fprintf(w, "Задача %s успешно удалена из БД!", TaskID)
}
