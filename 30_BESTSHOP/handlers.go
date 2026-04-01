package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type HttpHandler struct {
	Repo ShopRepository
}

func (h *HttpHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var u Users
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println("Не удалось прочитать запрос", err)
		return
	}

	userId, err := h.Repo.CreateUser(u.Name)
	if err != nil {
		log.Println("Не удалось создать пользователя", err)
		http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
		return
	}
	bankId, err := h.Repo.CreateBank(userId)
	if err != nil {
		log.Println("Не удалось создать банк", err)
		http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Пользователь создан с id %d и банк с id %d", userId, bankId)

}
func (h *HttpHandler) handleCreateCard(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Не удалось прочитать запрос", err)
		return
	}
	var bankID int
	var err error
	if req.UserID == 0 {
		bankID, err = h.Repo.CreateBank(0)
		if err != nil {
			log.Println("Не удалось создать банк", err)
			http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
			return
		}
	} else {
		bankID, err = h.Repo.GetBankByUserId(req.UserID)
		if err != nil {
			http.Error(w, "Не удалось найти банк", http.StatusInternalServerError)
			return
		}
	}
	cardNumber := fmt.Sprintf("%016d", rand.Int63n(1e16))
	cardID, err := h.Repo.CreateCard(bankID, cardNumber)
	if err != nil {
		log.Println("Не удалось создать карту", err)
		http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Карта создана с id %d", cardID)

}

func (h *HttpHandler) handleTopUp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CardID int    `json:"card_id"`
		Amount int    `json:"amount"`
		Valute string `json:"valute"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Не удалось прочитать запрос", err)
		http.Error(w, "Ошибка чтения JSON", http.StatusBadRequest)
		return
	}
	card, err := h.Repo.GetCard(req.CardID)
	if err != nil {
		log.Println("Не удалось найти карту", err)
		http.Error(w, "Карта не найдена", http.StatusNotFound)
		return
	}
	// 3. Выполняем бизнес-логику пополнения в зависимости от валюты
	if req.Valute == "RUB" {
		err = h.Repo.UpdateBalance(card.BankID, req.Amount, "RUB")
	} else if req.Valute == "USD" {
		amountUsd := req.Amount / 90 // конвертируем
		err = h.Repo.UpdateBalance(card.BankID, amountUsd, "USD")
	} else {
		http.Error(w, "Неизвестная валюта", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Println("Не удалось пополнить баланс", err)
		http.Error(w, "Ошибка при пополнении", http.StatusInternalServerError)
		return
	}
	// 4. Успешный ответ
	fmt.Fprintf(w, "Баланс карты успешно пополнен на %d %s", req.Amount, req.Valute)

}

func (h *HttpHandler) handleLinkCard(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int `json:"user_id"`
		CardID int `json:"card_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("Не удалось прочитать запрос", err)
		return
	}
	oldBank, err := h.Repo.GetBankById(req.CardID)
	if err != nil {
		log.Println("Не удалось найти банк", err)
		http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
		return
	}
	if oldBank.UserID != 0 {
		http.Error(w, "Карта уже привязана к другому пользователю!", http.StatusInternalServerError)
		return
	}
	userBankID, err := h.Repo.GetBankByUserId(req.UserID)
	if err != nil {
		http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
		return
	}
	err = h.Repo.UpdateCardBank(req.CardID, userBankID)
	if err != nil {
		http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
		return
	}
	if oldBank.BalanceRub > 0 {
		h.Repo.UpdateBalance(userBankID, oldBank.BalanceRub, "RUB")
		h.Repo.UpdateBalance(oldBank.ID, -oldBank.BalanceRub, "RUB")
	}
	if oldBank.BalanceUsd > 0 {
		h.Repo.UpdateBalance(userBankID, oldBank.BalanceUsd, "USD")
		h.Repo.UpdateBalance(oldBank.ID, -oldBank.BalanceUsd, "USD")
	}
	fmt.Fprintf(w, "Карта привязана к банку с id %d", userBankID)

}

func (h *HttpHandler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name   string `json:"name"`
		Price  int    `json:"price"`
		Amount int    `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Не удалось прочитать запрос", err)
		return
	}
	err := h.Repo.CreateProduct(req.Name, req.Price, req.Amount)
	if err != nil {
		log.Println("Не удалось создать продукт", err)
		http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Продукт создан")

}

func (h *HttpHandler) handleBuyProduct(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    int `json:"user_id"`
		ProductID int `json:"product_id"`
		Amount    int `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Не удалось прочитать запрос", err)
		return
	}
	product, err := h.Repo.GetProduct(req.ProductID)
	if err != nil {
		log.Println("Не удалось найти продукт", err)
		http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
		return
	}
	if product.Amount < req.Amount {
		http.Error(w, "Недостаточно товара на складе", http.StatusBadRequest)
		return
	}
	bankID, err := h.Repo.GetBankByUserId(req.UserID)
	if err != nil {
		http.Error(w, "Не удалось найти банк", http.StatusInternalServerError)
		return
	}
	balanceRub, err := h.Repo.GetBalance(bankID, "RUB")
	if err != nil {
		http.Error(w, "Не удалось найти баланс", http.StatusInternalServerError)
		return
	}
	if balanceRub < product.Price*req.Amount {
		http.Error(w, "Недостаточно средств на балансе", http.StatusBadRequest)
		return
	}
	err = h.Repo.UpdateBalance(bankID, -product.Price*req.Amount, "RUB")
	if err != nil {
		http.Error(w, "Не удалось обновить баланс", http.StatusInternalServerError)
		return
	}
	err = h.Repo.UpdateStock(req.ProductID, req.Amount)
	if err != nil {
		http.Error(w, "Не удалось обновить остаток", http.StatusInternalServerError)
		return
	}
	bonus := (product.Price * req.Amount) / 100
	err = h.Repo.UpdateBonusBalance(req.UserID, bonus)
	if err != nil {
		http.Error(w, "Не удалось начислить бонусы", http.StatusInternalServerError)
		return
	}
	err = h.Repo.CreateOrder(&Order{
		UserID:      req.UserID,
		ProductID:   req.ProductID,
		Amount:      req.Amount,
		Price:       product.Price * req.Amount,
		UsedBonuses: bonus,
		Status:      "completed",
	})
	if err != nil {
		http.Error(w, "Не удалось создать заказ", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Заказ успешно создан")

}

func (h *HttpHandler) autoUpdateOrders() {
	log.Println("Воркер автообновления заказов запущен...")
	for {
		time.Sleep(1 * time.Minute)

		orders, err := h.Repo.GetOrders()
		if err != nil {
			log.Println("Ошибка при получении заказов:", err)
			continue
		}
		for _, order := range orders {
			if order.Status == "paid" {
				err = h.Repo.UpdateOrderStatus(order.ID, "completed")
				if err != nil {
					log.Println("Ошибка при обновлении статуса заказа:", err)
					continue
				}
				log.Printf("Заказ %d успешно доставлен клиенту %d!", order.ID, order.UserID)
			}
		}
	}
}
func (h *HttpHandler) handleBuyWithBonus(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    int `json:"user_id"`
		ProductID int `json:"product_id"`
		Amount    int `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Не удалось прочитать запрос", err)
		return
	}
	product, err := h.Repo.GetProduct(req.ProductID)
	if err != nil {
		log.Println("Не удалось найти продукт", err)
		http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
		return
	}
	if product.Amount < req.Amount {
		http.Error(w, "Недостаточно товара на складе", http.StatusBadRequest)
		return
	}
	bankID, err := h.Repo.GetBankByUserId(req.UserID)
	if err != nil {
		http.Error(w, "Не удалось найти банк", http.StatusInternalServerError)
		return
	}
	balanceRub, err := h.Repo.GetBalance(bankID, "RUB")
	if err != nil {
		http.Error(w, "Не удалось найти баланс", http.StatusInternalServerError)
		return
	}
	bonusCard, err := h.Repo.GetBonusCard(req.UserID)
	if err != nil {
		http.Error(w, "Не удалось найти бонусную карту", http.StatusInternalServerError)
		return
	}
	if bonusCard.Status == "False" {
		http.Error(w, "У вас уже есть активный заказ с бонусами", http.StatusBadRequest)
		return
	}
	if bonusCard.Status == "True" {
		if bonusCard.BonusBalance > product.Price*req.Amount {
			bonusCard.BonusBalance = product.Price * req.Amount
		}
		needMoney := product.Price*req.Amount - bonusCard.BonusBalance
		if balanceRub < needMoney {
			http.Error(w, "Недостаточно средств на балансе", http.StatusBadRequest)
			return
		}
		err = h.Repo.UpdateBalance(bankID, -needMoney, "RUB")
		if err != nil {
			http.Error(w, "Не удалось обновить баланс", http.StatusInternalServerError)
			return
		}
		err = h.Repo.UpdateBonusBalance(req.UserID, -bonusCard.BonusBalance)
		if err != nil {
			http.Error(w, "Не удалось обновить баланс", http.StatusInternalServerError)
			return
		}
		err = h.Repo.SetBonusStatus(req.UserID, "False")
		if err != nil {
			http.Error(w, "Не удалось обновить статус карты", http.StatusInternalServerError)
			return
		}
	}
	err = h.Repo.UpdateStock(req.ProductID, req.Amount)
	if err != nil {
		http.Error(w, "Не удалось обновить остаток", http.StatusInternalServerError)
		return
	}
	bonus := (product.Price * req.Amount) / 100
	err = h.Repo.UpdateBonusBalance(req.UserID, bonus)
	if err != nil {
		http.Error(w, "Не удалось начислить бонусы", http.StatusInternalServerError)
		return
	}
	err = h.Repo.CreateOrder(&Order{
		UserID:      req.UserID,
		ProductID:   req.ProductID,
		Amount:      req.Amount,
		Price:       product.Price * req.Amount,
		UsedBonuses: bonusCard.BonusBalance,
		Status:      "paid",
	})
	if err != nil {
		http.Error(w, "Не удалось создать заказ", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Заказ успешно создан")

}
