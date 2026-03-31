package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=postgres dbname=BestShop host=localhost port=5432 sslmode=disable")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Ура к БД подсосались!")
}

type PostgresRepo struct {
	db *sql.DB
}

func (r *PostgresRepo) CreateUser(name string) (int, error) {
	var id int
	err := r.db.QueryRow("INSERT INTO users (name) VALUES ($1) RETURNING id", name).Scan(&id)
	return id, err
}

func (r *PostgresRepo) GetBankByUserId(userId int) (int, error) {
	var id int
	err := r.db.QueryRow("SELECT id FROM banks WHERE user_id = $1", userId).Scan(&id)
	return id, err
}

func (r *PostgresRepo) CreateBank(userId int) (int, error) {
	var id int
	if userId == 0 {
		err := r.db.QueryRow("INSERT INTO banks DEFAULT VALUES RETURNING id").Scan(&id)
		return id, err
	}
	err := r.db.QueryRow("INSERT INTO banks (user_id) VALUES ($1) RETURNING id", userId).Scan(&id)
	return id, err
}

func (r *PostgresRepo) UpdateBalance(bankId int, delta int, valute string) error {
	if valute == "RUB" {
		_, err := r.db.Exec("UPDATE banks SET balance_rub = balance_rub + $1 WHERE id = $2", delta, bankId)
		return err
	} else if valute == "USD" {
		_, err := r.db.Exec("UPDATE banks SET balance_usd = balance_usd + $1 WHERE id = $2", delta, bankId)
		return err
	}
	return nil
}

func (r *PostgresRepo) CreateCard(bankId int, number string) (int, error) {
	var id int
	err := r.db.QueryRow("INSERT INTO cards (bank_id, number) VALUES ($1, $2) RETURNING id", bankId, number).Scan(&id)
	return id, err
}

func (r *PostgresRepo) GetCard(cardId int) (*Cards, error) {
	var card Cards
	err := r.db.QueryRow("SELECT id, bank_id, number FROM cards WHERE id = $1", cardId).Scan(&card.ID, &card.BankID, &card.Number)
	return &card, err
}

func (r *PostgresRepo) GetProduct(id int) (*Products, error) {
	var product Products
	err := r.db.QueryRow("SELECT id, name, price, amount FROM products WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Price, &product.Amount)
	return &product, err
}

func (r *PostgresRepo) UpdateStock(productId int, delta int) error {
	_, err := r.db.Exec("UPDATE products SET amount = amount + $1 WHERE id = $2", delta, productId)
	return err
}

func (r *PostgresRepo) CreateOrder(order *Order) error {
	_, err := r.db.Exec("INSERT INTO orders (user_id, product_id, amount, price, used_bonuses, status) VALUES ($1, $2, $3, $4, $5, $6)", order.UserID, order.ProductID, order.Amount, order.Price, order.UsedBonuses, order.Status)
	return err
}

func (r *PostgresRepo) GetBonusCard(userId int) (*BonusCards, error) {
	var bonusCard BonusCards
	err := r.db.QueryRow("SELECT id, user_id, bonus_balance, status FROM bonuscards WHERE user_id = $1", userId).Scan(&bonusCard.ID, &bonusCard.UserID, &bonusCard.BonusBalance, &bonusCard.Status)
	return &bonusCard, err
}

func (r *PostgresRepo) UpdateBonusBalance(userId int, delta int) error {
	_, err := r.db.Exec("UPDATE bonuscards SET bonus_balance = bonus_balance + $1 WHERE user_id = $2", delta, userId)
	return err
}

func (r *PostgresRepo) SetBonusStatus(userId int, status string) error {
	_, err := r.db.Exec("UPDATE bonuscards SET status = $1 WHERE user_id = $2", status, userId)
	return err
}

func (r *PostgresRepo) GetOrder(userId int) (*Order, error) {
	var order Order
	err := r.db.QueryRow("SELECT id, user_id, product_id, amount, price, used_bonuses, status FROM orders WHERE user_id = $1", userId).Scan(&order.ID, &order.UserID, &order.ProductID, &order.Amount, &order.Price, &order.UsedBonuses, &order.Status)
	return &order, err
}

func (r *PostgresRepo) GetBalance(bankId int, valute string) (int, error) {
	var balance int
	var column string

	if valute == "RUB" {
		column = "balance_rub"
	} else {
		column = "balance_usd"
	}

	query := fmt.Sprintf("SELECT %s FROM banks WHERE id = $1", column)
	err := r.db.QueryRow(query, bankId).Scan(&balance)
	return balance, err
}
func (r *PostgresRepo) UpdateCardBank(cardId int, newBankId int) error {
	_, err := r.db.Exec("UPDATE cards SET bank_id = $1 WHERE id = $2", newBankId, cardId)
	return err
}
func (r *PostgresRepo) GetBankById(bankId int) (*Bank, error) {
	var bank Bank
	var nullUserId sql.NullInt32 // Специальная переменная для NULL-значений
	err := r.db.QueryRow("SELECT id, user_id, balance_rub, balance_usd FROM banks WHERE id = $1", bankId).Scan(&bank.ID, &nullUserId, &bank.BalanceRub, &bank.BalanceUsd)

	// Если user_id в базе не NULL, достаем настоящую цифру:
	if nullUserId.Valid {
		bank.UserID = int(nullUserId.Int32)
	} else {
		bank.UserID = 0 // Если NULL, ставим 0
	}

	return &bank, err
}
