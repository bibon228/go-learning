package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var u Users
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Println("Не удалось прочитать ID пользователя", err)
		return
	}
	err = db.QueryRow("INSERT INTO users (name) VALUES ($1) RETURNING id", u.Name).Scan(&u.ID)
	if err != nil {
		log.Println("Ошибка при создании пользователя", err)
		return
	}
	var bank Bank
	err = db.QueryRow("INSERT INTO banks (user_id) VALUES ($1) RETURNING id", u.ID).Scan(&bank.ID)
	if err != nil {
		log.Println("Не удалось создать аккаунт", err)
		return
	}
	_, err = db.Exec("INSERT INTO bonuscards (user_id, balance, status) VALUES ($1, 0, 'True')", u.ID)
	if err != nil {
		log.Println("Не удалось создать бонусную карту", err)

	}
	bank.UserID = u.ID
	bank.BalanceRub = 0
	bank.BalanceUsd = 0
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bank)

}
func handleCreateCard(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int `json:"user_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("Не удалось прочитать ID пользователя", err)
		return
	}
	var bankID int
	if request.UserID == 0 {
		err = db.QueryRow("INSERT INTO banks DEFAULT VALUES RETURNING id").Scan(&bankID)
		if err != nil {
			log.Println("Не удалось создать счет", err)
			return
		}
	} else {
		err = db.QueryRow("SELECT id FROM banks WHERE user_id = $1", request.UserID).Scan(&bankID)
		if err != nil {
			log.Println("Не удалось найти счет", err)
			return
		}
	}
	var card Cards
	cardNumber := fmt.Sprintf("%016d", rand.Int63n(1e16))
	err = db.QueryRow("INSERT INTO cards (bank_id, number) VALUES ($1, $2) RETURNING id", bankID, cardNumber).Scan(&card.ID)
	if err != nil {
		log.Println("Не удалось выдать карту", err)
		return
	}
	card.BankID = bankID
	card.Number = cardNumber

	json.NewEncoder(w).Encode(card)

}

func handleTopUp(w http.ResponseWriter, r *http.Request) {
	var request struct {
		CardID int    `json:"card_id"`
		Amount int    `json:"amount"`
		Valute string `json:"valute"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("Не удалось прочитать ID карты", err)
		return
	}
	var bankID int
	err = db.QueryRow("SELECT bank_id FROM cards WHERE id = $1", request.CardID).Scan(&bankID)
	if err != nil {
		log.Println("Не удалось найти аккаунт", err)
		return
	}
	if request.Valute == "RUB" {

		_, err = db.Exec("UPDATE banks SET balance_rub = balance_rub + $1 WHERE id = $2", request.Amount, bankID)
		if err != nil {
			log.Println("Не удалось обновить баланс", err)
			return
		}
	} else if request.Valute == "USD" {
		amount := request.Amount / 90
		_, err = db.Exec("UPDATE banks SET balance_usd = balance_usd + $1 WHERE id = $2", amount, bankID)
		if err != nil {
			log.Println("Не удалось обновить баланс", err)
			return
		}
	}

	fmt.Fprintf(w, "Баланс карты %d успешно пополнен на %d\n", request.CardID, request.Amount)
}
func handleLinkCard(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int `json:"user_id"`
		CardID int `json:"card_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("Не удалось прочитать ID пользователя", err)
		return
	}
	var userBankID int
	err = db.QueryRow("SELECT id FROM banks WHERE user_id = $1", request.UserID).Scan(&userBankID)
	if err != nil {
		log.Println("Не удалось найти банковский счет", err)
		return
	}
	var cardBankID int
	err = db.QueryRow("SELECT bank_id FROM cards WHERE id = $1", request.CardID).Scan(&cardBankID)
	if err != nil {
		log.Println("Не удалось найти банковский счет карты", err)
		return
	}
	var cardUserID *int
	err = db.QueryRow("SELECT user_id FROM banks WHERE id = $1", cardBankID).Scan(&cardUserID)
	if err != nil {
		log.Println("Не удалось найти пользователя", err)
		return
	}
	if cardUserID != nil {
		fmt.Fprintf(w, "Карта %d уже привязана к другому пользователю", request.CardID)
		return
	}
	err = db.QueryRow("UPDATE cards SET bank_id = $1 WHERE id = $2 RETURNING id", userBankID, request.CardID).Scan(&request.CardID)
	if err != nil {
		log.Println("Не удалось привязать карту", err)
		return
	}
	fmt.Fprintf(w, "Карта %d успешно привязана к пользователю %d", request.CardID, request.UserID)

	var balance_rub int
	var balance_usd int
	err = db.QueryRow("SELECT balance_rub, balance_usd FROM banks WHERE id = $1", cardBankID).Scan(&balance_rub, &balance_usd)
	if err != nil {
		log.Println("Не удалось найти баланс", err)
		return
	}
	_, err = db.Exec("UPDATE banks SET balance_rub = balance_rub + $1, balance_usd = balance_usd + $2 WHERE id = $3", balance_rub, balance_usd, userBankID)
	if err != nil {
		log.Println("Не удалось обновить баланс", err)
		return
	}
	_, err = db.Exec("UPDATE banks SET balance_rub = 0, balance_usd = 0 WHERE id = $1", cardBankID)
	if err != nil {
		log.Println("Не удалось обнулить баланс", err)
		return
	}
}

func handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var p Products
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := db.Exec("INSERT INTO products (name, price, amount) VALUES ($1, $2, $3)", p.Name, p.Price, p.Amount)
	if err != nil {
		log.Println("Ошибка при создании продукта", err)
		return
	}
	fmt.Fprintf(w, "Продукт %s успешно создан!", p.Name)
}
func handleBuyProduct(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID    int `json:"user_id"`
		ProductID int `json:"product_id"`
		Amount    int `json:"amount"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("Не удалось прочитать ID пользователя", err)
		return
	}
	var product Products
	err = db.QueryRow("SELECT * FROM products WHERE id = $1", request.ProductID).Scan(&product.ID, &product.Name, &product.Price, &product.Amount)
	if err != nil {
		log.Println("Не удалось найти продукт", err)
		return
	}
	if product.Amount < request.Amount {
		fmt.Fprintf(w, "Недостаточно товара на складе")
		return
	}
	var userBankID int
	err = db.QueryRow("SELECT id FROM banks WHERE user_id = $1", request.UserID).Scan(&userBankID)
	if err != nil {
		log.Println("Не удалось найти банковский счет", err)
		return
	}
	var balance_rub int
	err = db.QueryRow("SELECT balance_rub FROM banks WHERE id = $1", userBankID).Scan(&balance_rub)
	if err != nil {
		log.Println("Не удалось найти баланс", err)
		return
	}
	if balance_rub < product.Price*request.Amount {
		fmt.Fprintf(w, "Недостаточно средств на балансе")
		return
	}
	_, err = db.Exec("UPDATE banks SET balance_rub = balance_rub - $1 WHERE id = $2", product.Price*request.Amount, userBankID)
	if err != nil {
		log.Println("Не удалось обновить баланс", err)
		return
	}
	_, err = db.Exec("UPDATE products SET amount = amount - $1 WHERE id = $2", request.Amount, request.ProductID)
	if err != nil {
		log.Println("Не удалось обновить количество товара", err)
		return
	}
	var bonus int
	if product.Price*request.Amount < 1000 {
		bonus = product.Price * request.Amount / 200
		_, err = db.Exec("UPDATE bonuscards SET balance = balance + $1 WHERE user_id = $2", bonus, request.UserID)
	} else {
		bonus = product.Price * request.Amount / 100
		_, err = db.Exec("UPDATE bonuscards SET balance = balance + $1 WHERE user_id = $2", bonus, request.UserID)
	}
	if err != nil {
		log.Println("Не удалось обновить баланс", err)
		return
	}
	var order Order
	order.UserID = request.UserID
	order.ProductID = request.ProductID
	order.Amount = request.Amount
	order.Price = product.Price * request.Amount
	order.UsedBonuses = 0
	order.Status = "paid"
	_, err = db.Exec("INSERT INTO orders (user_id, product_id, amount, price, used_bonuses, status) VALUES ($1, $2, $3, $4, $5, $6)", order.UserID, order.ProductID, order.Amount, order.Price, order.UsedBonuses, order.Status)
	if err != nil {
		log.Println("Не удалось создать заказ", err)
		return
	}
	fmt.Fprintf(w, "Товар %s успешно куплен!", product.Name)

}

func autoUpdateOrders() {
	var order Order
	for {
		time.Sleep(1 * time.Minute)
		rows, err := db.Query("SELECT id, user_id, product_id, amount, price, used_bonuses, status FROM orders WHERE status = 'paid'")
		if err != nil {
			log.Println("Не удалось найти заказы", err)
			continue
		}
		for rows.Next() {
			err = rows.Scan(&order.ID, &order.UserID, &order.ProductID, &order.Amount, &order.Price, &order.UsedBonuses, &order.Status)
			if err != nil {
				log.Println("Не удалось прочитать заказ", err)
				continue

			}
			if order.Status == "paid" {
				order.Status = "completed"
				_, err := db.Exec("UPDATE orders SET status = $1 WHERE id = $2", order.Status, order.ID)
				if err != nil {
					log.Println("Не удалось обновить статус заказа", err)
					return
				}
				_, err = db.Exec("UPDATE bonuscards SET status = 'True' WHERE user_id = $1", order.UserID)
				if err != nil {
					log.Println("Не удалось обновить статус", err)
					return
				}
				_, err = db.Exec("UPDATE bonuscards SET balance = balance - $1 WHERE user_id = $2", order.UsedBonuses, order.UserID)
				if err != nil {
					log.Println("Не удалось обновить баланс", err)
					return
				}
			}
		}
		rows.Close()
	}

}
func handleBuyWithBonus(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID    int `json:"user_id"`
		ProductID int `json:"product_id"`
		Amount    int `json:"amount"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("Не удалось прочитать ID пользователя", err)
		return
	}
	var product Products
	err = db.QueryRow("SELECT * FROM products WHERE id = $1", request.ProductID).Scan(&product.ID, &product.Name, &product.Price, &product.Amount)
	if err != nil {
		log.Println("Не удалось найти продукт", err)
		return
	}
	if product.Amount < request.Amount {
		fmt.Fprintf(w, "Недостаточно товара на складе")
		return
	}
	var userBankID int
	err = db.QueryRow("SELECT id FROM banks WHERE user_id = $1", request.UserID).Scan(&userBankID)
	if err != nil {
		log.Println("Не удалось найти банковский счет", err)
		return
	}
	var balance_rub int
	err = db.QueryRow("SELECT balance_rub FROM banks WHERE id = $1", userBankID).Scan(&balance_rub)
	if err != nil {
		log.Println("Не удалось найти баланс", err)
		return
	}
	var bonusBalance int
	err = db.QueryRow("SELECT balance FROM bonuscards WHERE user_id = $1", request.UserID).Scan(&bonusBalance)
	if err != nil {
		log.Println("Не удалось найти баланс", err)
		return
	}

	var bonusStatus string
	err = db.QueryRow("SELECT status FROM bonuscards WHERE user_id = $1", request.UserID).Scan(&bonusStatus)
	if bonusStatus == "False" {
		fmt.Fprintf(w, "У вас уже есть активный заказ с бонусами")
		return
	}
	if bonusStatus == "True" {
		if bonusBalance > product.Price*request.Amount {
			bonusBalance = product.Price * request.Amount
		}
		needMoney := product.Price*request.Amount - bonusBalance
		if balance_rub < needMoney {
			fmt.Fprintf(w, "Недостаточно средств на балансе")
			return
		}

		balance := balance_rub - needMoney
		_, err = db.Exec("UPDATE banks SET balance_rub = $1 WHERE id = $2", balance, userBankID)
		if err != nil {
			log.Println("Не удалось обновить баланс", err)
			return
		}
		_, err = db.Exec("UPDATE products SET amount = amount - $1 WHERE id = $2", request.Amount, request.ProductID)
		if err != nil {
			log.Println("Не удалось обновить количество товара", err)
			return
		}
		var bonus int
		if product.Price*request.Amount < 1000 {
			bonus = product.Price * request.Amount / 200
			_, err = db.Exec("UPDATE bonuscards SET balance = balance + $1 WHERE user_id = $2", bonus, request.UserID)
		} else {
			bonus = product.Price * request.Amount / 100
			_, err = db.Exec("UPDATE bonuscards SET balance = balance + $1 WHERE user_id = $2", bonus, request.UserID)
		}
		if err != nil {
			log.Println("Не удалось обновить баланс", err)
			return
		}
		var order Order
		order.UserID = request.UserID
		order.ProductID = request.ProductID
		order.Amount = request.Amount
		order.Price = product.Price * request.Amount
		order.UsedBonuses = bonusBalance
		order.Status = "paid"
		_, err = db.Exec("INSERT INTO orders (user_id, product_id, amount, price, used_bonuses, status) VALUES ($1, $2, $3, $4, $5, $6)", order.UserID, order.ProductID, order.Amount, order.Price, order.UsedBonuses, order.Status)
		if err != nil {
			log.Println("Не удалось создать заказ", err)
			return
		}
		_, err = db.Exec("UPDATE bonuscards SET status = 'False' WHERE user_id = $1", request.UserID)
		if err != nil {
			log.Println("Ошибка при обновлении статуса бонусов", err)
			return
		}

		fmt.Fprintf(w, "Товар %s успешно куплен!", product.Name)

	}

}
