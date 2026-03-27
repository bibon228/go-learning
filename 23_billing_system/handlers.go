package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal("Не удалось прочитать пользователя", err)
		return
	}

	err = DB.QueryRow("INSERT INTO users (name) VALUES ($1) RETURNING id", user.Name).Scan(&user.ID)
	if err != nil {
		log.Fatal("Не удалось создать пользователя", err)
		return
	}
	_, err = DB.Exec("INSERT INTO accounts (user_id) VALUES ($1)", user.ID)
	if err != nil {
		log.Fatal("Не удалось создать аккаунт", err)
		return
	}
	fmt.Fprintf(w, "Пользователь %s успешно создан! Его ID: %d\n", user.Name, user.ID)

}
func handleIssueCard(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int `json:"user_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal("Не удалось прочитать ID пользователя", err)
		return
	}
	var accountID int
	if request.UserID == 0 {
		err = DB.QueryRow("INSERT INTO accounts DEFAULT VALUES RETURNING id").Scan(&accountID)
		if err != nil {
			log.Fatal("Не удалось создать аккаунт", err)
			return
		}
	} else {
		err = DB.QueryRow("SELECT id FROM accounts WHERE user_id = $1", request.UserID).Scan(&accountID)
		if err != nil {
			log.Fatal("Не удалось найти аккаунт", err)
			return
		}
	}
	var card Card
	err = DB.QueryRow("INSERT INTO cards (account_id) VALUES ($1) RETURNING id", accountID).Scan(&card.ID)
	if err != nil {
		log.Fatal("Не удалось выдать карту", err)
		return
	}
	card.AccountID = accountID

	json.NewEncoder(w).Encode(card)

}

func handleTopUp(w http.ResponseWriter, r *http.Request) {
	var request struct {
		CardID int `json:"card_id"`
		Amount int `json:"amount"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal("Не удалось прочитать ID карты", err)
		return
	}
	var account_id int
	err = DB.QueryRow("SELECT account_id FROM cards WHERE id = $1", request.CardID).Scan(&account_id)
	if err != nil {
		log.Fatal("Не удалось найти аккаунт", err)
		return
	}

	_, err = DB.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", request.Amount, account_id)
	if err != nil {
		log.Fatal("Не удалось обновить баланс", err)
		return
	}
	var transaction Transaction
	transaction.AccountID = account_id
	transaction.Amount = request.Amount
	transaction.Type = "top_up"
	_, err = DB.Exec("INSERT INTO transactions (account_id, amount, type) VALUES ($1, $2, $3)", transaction.AccountID, transaction.Amount, transaction.Type)
	if err != nil {
		log.Fatal("Не удалось создать транзакцию", err)
		return
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
		log.Fatal("Не удалось прочитать ID пользователя", err)
		return
	}
	var userAccountID int
	err = DB.QueryRow("SELECT id FROM accounts WHERE user_id = $1", request.UserID).Scan(&userAccountID)
	if err != nil {
		log.Fatal("Не удалось найти аккаунт", err)
		return
	}
	var cardAccountID int
	err = DB.QueryRow("SELECT account_id FROM cards WHERE id = $1", request.CardID).Scan(&cardAccountID)
	if err != nil {
		log.Fatal("Не удалось найти аккаунт", err)
		return
	}
	var cardUserID *int
	err = DB.QueryRow("SELECT user_id FROM accounts WHERE id = $1", cardAccountID).Scan(&cardUserID)
	if err != nil {
		log.Fatal("Не удалось найти пользователя", err)
		return
	}
	if cardUserID != nil {
		fmt.Fprintf(w, "Карта %d уже привязана к другому пользователю", request.CardID)
		return
	}
	err = DB.QueryRow("UPDATE cards SET account_id = $1 WHERE id = $2 RETURNING id", userAccountID, request.CardID).Scan(&request.CardID)
	if err != nil {
		log.Fatal("Не удалось привязать карту", err)
		return
	}
	fmt.Fprintf(w, "Карта %d успешно привязана к пользователю %d", request.CardID, request.UserID)

	var balance int
	err = DB.QueryRow("SELECT balance FROM accounts WHERE id = $1", cardAccountID).Scan(&balance)
	if err != nil {
		log.Fatal("Не удалось найти баланс", err)
		return
	}
	_, err = DB.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", balance, userAccountID)
	if err != nil {
		log.Fatal("Не удалось обновить баланс", err)
		return
	}
	_, err = DB.Exec("UPDATE accounts SET balance = 0 WHERE id = $1", cardAccountID)
	if err != nil {
		log.Fatal("Не удалось обнулить баланс", err)
		return
	}
	fmt.Fprintf(w, "Баланс карты %d успешно перенесен на баланс пользователя %d", request.CardID, request.UserID)
	var transaction Transaction
	transaction.AccountID = userAccountID
	transaction.Amount = balance
	transaction.Type = "transfer"
	_, err = DB.Exec("INSERT INTO transactions (account_id, amount, type) VALUES ($1, $2, $3)", transaction.AccountID, transaction.Amount, transaction.Type)
	if err != nil {
		log.Fatal("Не удалось создать транзакцию", err)
		return
	}
	fmt.Fprintf(w, "Баланс карты %d успешно перенесен на баланс пользователя %d\n", request.CardID, request.UserID)

}
