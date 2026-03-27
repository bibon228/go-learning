package main

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Done   bool   `json:"done"`
	UserID int    `json:"user_id"`
}

type CreateUserRequest struct {
	Name string `json:"name"`
}
type CreateTaskRequest struct {
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
}
type UpdateTaskStatusRequest struct {
	ID     int  `json:"id"`
	TaskID int  `json:"task_id"`
	Done   bool `json:"done"`
}
