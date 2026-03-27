package main

import "fmt"

type Blacklist map[string]bool

func IsAllowed(name string, list Blacklist) bool {
	_, ok := list[name]
	if ok {
		return false
	}
	return true
}

type Product struct {
	Name  string
	Price int
}

func ApplyDiscount(p *Product, amount int) {
	p.Price -= amount
}

func ChangeName(p *Product, newName string) {
	p.Name = newName
}

func main() {
	// Работа с мапой
	myList := Blacklist{
		"Gimli": true,
		"Orc":   true,
	}
	fmt.Println("Можно ли войти Gimli? :", IsAllowed("Gimli", myList))
	fmt.Println("Можно ли войти LiveCoder? :", IsAllowed("LiveCoder", myList))

	laptop := Product{Name: "Laptop", Price: 1000} // Теперь имя совпадает с type Product

	ChangeName(&laptop, "Smartphone")
	fmt.Println("Новое имя:", laptop.Name)

	ApplyDiscount(&laptop, 100)
	fmt.Println("Новая цена:", laptop.Price)
}
