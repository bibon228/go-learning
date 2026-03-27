package main

import "fmt"

// Твой код для Зоопарка пойдет сюда

type Entity interface {
	Add() string
	Remove() string
	Update() string
	List() string
}

type Task struct {
	ID     int
	Title  string
	Done   bool
	UserID int
}

type Users struct {
	ID   int
	Name string
}

func (t *Task) Add() string {
	return "Задача добавлена"
}

func (t *Task) Remove() string {
	return "Задача удалена"
}

func (t *Task) Update() string {
	return "Задача обновлена"
}

func (t *Task) List() string {
	return "Список задач"
}
func (u *Users) Add() string {
	return "Пользователь добавлен"
}
func (u *Users) Remove() string {
	return "Пользователь удален"
}
func (u *Users) Update() string {
	return "Пользователь обновлен"
}
func (u *Users) List() string {
	return "Список пользователей"
}

func PrintTasks(e Entity) {
	fmt.Println(e.Add())
	fmt.Println(e.Remove())
	fmt.Println(e.Update())
	fmt.Println(e.List())
}
func main() {
	user := &Users{
		ID:   1,
		Name: "Пользователь 1",
	}
	PrintTasks(user)
}
