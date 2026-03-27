package main

import (
	"basics/models"
	"fmt"
)

func main() {
	fmt.Println("База студентов, доступные команды : add, remove, print, infoStudent, exit")
	students := make(map[string]*models.Student)
	for {
		var command string
		fmt.Scan(&command)
		switch command {
		case "add":
			var name string
			var age int
			fmt.Scan(&name, &age)
			students[name] = models.NewStudent(name, age)
			fmt.Println("Студент добавлен")
		case "remove":
			var name string
			fmt.Scan(&name)
			delete(students, name)
			fmt.Println("Студент удален")
		case "print":
			if len(students) == 0 {
				fmt.Println("Список студентов пуст")
			} else {
				for name, student := range students {
					fmt.Println("Имя: ", name, "Возраст: ", student.Age)
				}
			}

		case "infoStudent":
			var name string
			fmt.Scan(&name)
			info := students[name]
			if info != nil {
				fmt.Println("Имя: ", info.Name, "Возраст: ", info.Age)
			} else {
				fmt.Println("Студент не найден")
			}
		case "exit":
			return
		}
	}

	// Тут будет твой бесконечный цикл
}
