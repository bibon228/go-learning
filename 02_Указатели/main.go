package main

import "fmt"

type Blacklist map[string]bool

func IsAllowed(name string, list Blacklist) bool {
	_, ok := list[name]
	if ok {
		return false
	}
	return true

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
	myList := Blacklist{
		"Gimli": true,
		"Orc":   true,
	}
	yesOrNo := IsAllowed("Gimli", myList)
	fmt.Println("Можно ли войти Gimli? :", yesOrNo)
	yesOrNoCoder := IsAllowed("LiveCoder", myList)
	fmt.Println("Можно ли войти LiveCoder? :", yesOrNoCoder)
	laptop := Product{Name: "Laptop", Price: 1000}

	ChangeName(&laptop, "Smartphone")
	fmt.Println("Новое имя:", laptop.Name)
	ApplyDiscount(&laptop, 100)
	fmt.Println("Новая цена:", laptop.Price)
}