package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handleAddUser(w http.ResponseWriter, r *http.Request) {
	var user Users
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal("Не удалось прочитать пользователя", err)
		return
	}
	err = DB.QueryRow("INSERT INTO users (name) VALUES ($1) RETURNING id ", user.Name).Scan(&user.ID)
	if err != nil {
		log.Fatal("Ошибка при добавление пользователя", err)
		return
	}
	_, err = DB.Exec("INSERT INTO bank (user_id) VALUES ($1)", user.ID)
	if err != nil {
		log.Fatal("Ошибка при создании счета пользователю", err)
		return

	}
	fmt.Fprintf(w, "Пользователь %s успешно создан! Его ID: %d\n", user.Name, user.ID)

}
func handleCreateCard(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int `json:"user_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal("Ошиибка при создании карты пользователю", err)
		return
	}
	var bankID int
	if request.UserID == 0 {
		err = DB.QueryRow("INSERT INTO bank DEFAULT VALUES RETURNING id").Scan(&bankID)
		if err != nil {
			log.Fatal("Не удалось создать анонимный счет!", err)
			return
		}
	} else {
		err = DB.QueryRow("SELECT id FROM bank WHERE user_id = $1", request.UserID).Scan(&bankID)
		if err != nil {
			log.Fatal("Не удалось найти счет данного пользователя!", err)
			return
		}
	}

	var card Cards
	err = DB.QueryRow("INSERT INTO cards (bank_id) VALUES ($1) RETURNING id", bankID).Scan(&card.ID)
	if err != nil {
		log.Fatal("Не удалось создать карту!", err)
		return
	}
	card.BankID = bankID
	json.NewEncoder(w).Encode(card)

}

func handleToPup(w http.ResponseWriter, r *http.Request) {
	var request struct {
		CardID int `json:"card_id"`
		Amount int `json:"amount"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal("Ошибка при вводе данных", err)
		return
	}

	var bankID int
	err = DB.QueryRow("SELECT bank_id FROM cards WHERE id = $1", request.CardID).Scan(&bankID)
	if err != nil {
		log.Fatal("Ошибка при получении счета", err)
		return
	}
	_, err = DB.Exec("UPDATE bank SET balance = balance + $1 WHERE id = $2", request.Amount, bankID)
	if err != nil {
		log.Fatal("Ошибка при пополнение счета карты!", err)
		return
	}
	var logs Logs
	logs.BankID = bankID
	logs.Amount = request.Amount
	logs.Info = "Пополнение"
	_, err = DB.Exec("INSERT INTO logs (bank_id, amount, info) VALUES ($1, $2, $3)", logs.BankID, logs.Amount, logs.Info)
	if err != nil {
		log.Fatal("Ошибка при логировании!", err)
		return
	}
	fmt.Fprintf(w, "Карта с ID %d успешно пополненна на %d\n", request.CardID, request.Amount)

}
func handleLinkCard(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID int `json:"user_id"`
		CardID int `json:"card_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal("Ошибка при вводе данных")
		return
	}
	var userBankID int
	err = DB.QueryRow("SELECT id FROM bank WHERE user_id = $1 ", request.UserID).Scan(&userBankID)
	if err != nil {
		log.Fatal("Не удалось найти счет пользователя!", err)
		return
	}
	var cardBankID int
	err = DB.QueryRow("SELECT bank_id FROM cards WHERE id = $1 ", request.CardID).Scan(&cardBankID)
	if err != nil {
		log.Fatal("Ошибка при получении счета карты!", err)
		return
	}
	var cardUserID *int
	err = DB.QueryRow("SELECT user_id FROM bank WHERE id = $1", cardBankID).Scan(&cardUserID)
	if err != nil {
		log.Fatal("Ошибка при получении информации о владельце карты!", err)
		return

	}
	if cardUserID != nil {
		fmt.Fprintf(w, "Карта уже привязана к пользователю: %d\n", cardUserID)
		return

	}
	var balance int
	err = DB.QueryRow("SELECT balance FROM bank WHERE id = $1", cardBankID).Scan(&balance)
	if err != nil {
		log.Fatal("Ошибка при получении данных о балансе карты!", err)
		return
	}
	_, err = DB.Exec("UPDATE bank SET balance = balance + $1 WHERE id = $2", balance, userBankID)
	if err != nil {
		log.Fatal("Ошибка при переводе денег с бонусной карты!", err)
		return
	}
	_, err = DB.Exec("UPDATE bank SET balance = 0 WHERE id = $1", cardBankID)
	if err != nil {
		log.Fatal("Ошибка при обнулении бонусной карты!", err)
		return
	}
	var logs Logs
	logs.BankID = userBankID
	logs.Amount = balance
	logs.Info = "Активация бонусной карты"
	_, err = DB.Exec("INSERT INTO logs (bank_id, amount, info) VALUES ($1, $2, $3)", logs.BankID, logs.Amount, logs.Info)
	if err != nil {
		log.Fatal("Ошибка при логировании!", err)
		return
	}
	fmt.Fprintf(w, "Баланс с карты %d успешно переведен на счет пользователя %d\n", request.CardID, request.UserID)

}
