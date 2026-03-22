package main

import "fmt"

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
	laptop := Product{Name: "Laptop", Price: 1000}

	ChangeName(&laptop, "Smartphone")
	fmt.Println("Новое имя:", laptop.Name)
	ApplyDiscount(&laptop, 100)
	fmt.Println("Новая цена:", laptop.Price)
}
