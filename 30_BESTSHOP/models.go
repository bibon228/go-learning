package main

type Users struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Bank struct {
	ID         int `json:"id"`
	UserID     int `json:"user_id"`
	BalanceRub int `json:"balance_rub"`
	BalanceUsd int `json:"balance_usd"`
}
type Cards struct {
	ID      int    `json:"id"`
	BankID  int    `json:"bank_id"`
	Number  string `json:"number"`
	Balance int    `json:"balance"`
}
type Products struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Amount int    `json:"amount"`
}
type Order struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	ProductID   int    `json:"product_id"`
	Amount      int    `json:"amount"`
	Price       int    `json:"price"`
	UsedBonuses int    `json:"used_bonuses"`
	Status      string `json:"status"`
}
type BonusCards struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	BonusBalance int    `json:"bonus_balance"`
	Status       string `json:"status"`
}

type ShopRepository interface {
	// РАБОТА С ЮЗЕРАМИ И БАНКОМ
	CreateUser(name string) (int, error)
	GetBankByUserId(userId int) (int, error)
	GetBalance(bankId int, valute string) (int, error)
	UpdateBalance(bankId int, delta int, valute string) error
	CreateBank(userId int) (int, error)
	GetBankById(bankId int) (*Bank, error)
	UpdateCardBank(cardId int, newBankId int) error

	// КАРТЫ
	CreateCard(bankId int, number string) (int, error)
	GetCard(cardId int) (*Cards, error)

	// ТОВАРЫ И ЗАКАЗЫ
	GetProduct(id int) (*Products, error)
	UpdateStock(productId int, delta int) error
	CreateOrder(order *Order) error
	GetOrder(userId int) (*Order, error)
	GetOrders() ([]Order, error)
	UpdateOrderStatus(orderId int, status string) error
	CreateProduct(name string, price int, amount int) error
	// БОНУСЫ
	GetBonusCard(userId int) (*BonusCards, error)
	UpdateBonusBalance(userId int, delta int) error
	SetBonusStatus(userId int, status string) error
}
