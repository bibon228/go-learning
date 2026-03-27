package main

import "fmt"

// Твоему коду нужно появиться здесь:
// 1. Интерфейс PaymentMethod
// 2. Структура CreditCard
// 3. Структура PayPal
// 4. Метод Pay(amount int) для каждой из них

// 5. Общая функция ProcessPayment
// func ProcessPayment(...) {
//	 ...
// }

type PaymentMethod interface {
	Pay(amount int) string
}

type CreditCard struct {
	Number string

	Expiry string
}

func (c *CreditCard) Pay(amount int) string {
	return fmt.Sprintf("Оплата картой %s на сумму %d", c.Number, amount)
}

type PayPal struct {
	Email    string
	Password string
}

func (p *PayPal) Pay(amount int) string {
	return fmt.Sprintf("Оплата PayPal %s на сумму %d", p.Email, amount)
}
func ProcessPayment(p PaymentMethod, amount int) {
	result := p.Pay(amount)
	fmt.Println(result)
}

func main() {
	// 6. Создаем карту и пейпал, скармливаем их функции ProcessPayment
	card := &CreditCard{
		Number: "1234567890123456",

		Expiry: "12/25",
	}
	paypal := &PayPal{
		Email:    "[EMAIL_ADDRESS]",
		Password: "password",
	}
	ProcessPayment(card, 100)
	ProcessPayment(paypal, 200)

}
