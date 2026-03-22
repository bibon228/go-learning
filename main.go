package main

import (
	"fmt"
	"net/http"
)

type Product struct {
	Name  string
	Price int
}

var laptop = Product{Name: "MacBook", Price: 2000}

func applyDiscount(p *Product, amount int) {
	p.Price -= amount
}
func handleDiscount(w http.ResponseWriter, r *http.Request) {
	applyDiscount(&laptop, 100)
	fmt.Fprintf(w, "Текущая цена: %d\n", laptop.Price)
}
func main() {
	http.HandleFunc("/discount", handleDiscount)
	http.ListenAndServe(":8080", nil)

}
