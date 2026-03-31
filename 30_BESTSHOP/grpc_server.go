package main

import (
	"context"
	shop "db/proto"
	"fmt"
	"log"
	"math/rand"
)

type ShopServer struct {
	shop.UnimplementedShopServiceServer
	Repo ShopRepository
}

func (s *ShopServer) BuyProduct(ctx context.Context, req *shop.BuyRequest) (*shop.BuyResponse, error) {
	log.Printf("Пришёл запрос на покупку: UserID=%d", req.UserId)
	product, err := s.Repo.GetProduct(int(req.ProductId))
	if err != nil {
		return &shop.BuyResponse{Success: false, Message: "Продукт не найден"}, nil
	}
	if product.Amount < int(req.Amount) {
		return &shop.BuyResponse{Success: false, Message: "Недостаточно товара на складе"}, nil
	}
	bankID, err := s.Repo.GetBankByUserId(int(req.UserId))
	if err != nil {
		return &shop.BuyResponse{Success: false, Message: "Не удалось найти счет"}, nil
	}
	balance_rub, err := s.Repo.GetBalance(bankID, "RUB")
	if err != nil {
		return &shop.BuyResponse{Success: false, Message: "Не удалось найти баланс"}, nil
	}
	if balance_rub < product.Price*int(req.Amount) {
		return &shop.BuyResponse{Success: false, Message: "Недостаточно средств на балансе"}, nil
	}
	err = s.Repo.UpdateBalance(bankID, product.Price*int(req.Amount), "RUB")
	if err != nil {
		return &shop.BuyResponse{Success: false, Message: "Не удалось обновить баланс"}, nil
	}
	err = s.Repo.UpdateStock(int(req.ProductId), int(req.Amount))
	if err != nil {
		return &shop.BuyResponse{Success: false, Message: "Не удалось обновить остаток"}, nil
	}
	bonus := (product.Price * int(req.Amount)) / 100
	err = s.Repo.UpdateBonusBalance(int(req.UserId), bonus)
	if err != nil {
		return &shop.BuyResponse{Success: false, Message: "Не удалось начислить бонусы"}, nil
	}
	err = s.Repo.CreateOrder(&Order{
		UserID:      int(req.UserId),
		ProductID:   int(req.ProductId),
		Amount:      int(req.Amount),
		Price:       product.Price * int(req.Amount),
		UsedBonuses: bonus,
		Status:      "completed",
	})
	if err != nil {
		return &shop.BuyResponse{Success: false, Message: "Не удалось создать заказ"}, nil
	}

	return &shop.BuyResponse{
		Success: true,
		Message: "Всё сработало!",
	}, nil
}
func (s *ShopServer) BuyWithBonus(ctx context.Context, req *shop.BuyWithBonusRequest) (*shop.BuyWithBonusResponse, error) {
	log.Printf("Пришёл запрос на покупку: UserID=%d", req.UserId)
	product, err := s.Repo.GetProduct(int(req.ProductId))
	if err != nil {
		return &shop.BuyWithBonusResponse{Success: false, Message: "Продукт не найден"}, nil
	}
	if product.Amount < int(req.Amount) {
		return &shop.BuyWithBonusResponse{Success: false, Message: "Недостаточно товара на складе"}, nil
	}
	bankID, err := s.Repo.GetBankByUserId(int(req.UserId))
	if err != nil {
		return &shop.BuyWithBonusResponse{Success: false, Message: "Не удалось найти счет"}, nil
	}
	balance_rub, err := s.Repo.GetBalance(bankID, "RUB")
	if err != nil {
		return &shop.BuyWithBonusResponse{Success: false, Message: "Не удалось найти баланс"}, nil
	}
	bonusCard, err := s.Repo.GetBonusCard(int(req.UserId))
	if bonusCard.Status == "False" {
		return &shop.BuyWithBonusResponse{Success: false, Message: "У вас уже есть активный заказ с бонусами"}, nil
	}
	if bonusCard.Status == "True" {
		if bonusCard.BonusBalance > product.Price*int(req.Amount) {
			bonusCard.BonusBalance = product.Price * int(req.Amount)
		}
		needMoney := product.Price*int(req.Amount) - bonusCard.BonusBalance
		if balance_rub < needMoney {
			return &shop.BuyWithBonusResponse{Success: false, Message: "Недостаточно средств на балансе"}, nil
		}
		err = s.Repo.UpdateBalance(bankID, -needMoney, "RUB")
		if err != nil {
			return &shop.BuyWithBonusResponse{Success: false, Message: "Не удалось обновить баланс"}, nil
		}
		err = s.Repo.UpdateBonusBalance(int(req.UserId), -bonusCard.BonusBalance)
		if err != nil {
			return &shop.BuyWithBonusResponse{Success: false, Message: "Не удалось обновить баланс"}, nil
		}
		err = s.Repo.SetBonusStatus(int(req.UserId), "False")
		if err != nil {
			return &shop.BuyWithBonusResponse{Success: false, Message: "Не удалось обновить статус карты"}, nil
		}
	}
	err = s.Repo.UpdateStock(int(req.ProductId), int(req.Amount))
	if err != nil {
		return &shop.BuyWithBonusResponse{Success: false, Message: "Не удалось обновить остаток"}, nil
	}
	bonus := (product.Price * int(req.Amount)) / 100
	err = s.Repo.UpdateBonusBalance(int(req.UserId), bonus)
	if err != nil {
		return &shop.BuyWithBonusResponse{Success: false, Message: "Не удалось начислить бонусы"}, nil
	}
	s.Repo.CreateOrder(&Order{
		UserID:      int(req.UserId),
		ProductID:   int(req.ProductId),
		Amount:      int(req.Amount),
		Price:       product.Price * int(req.Amount),
		UsedBonuses: bonusCard.BonusBalance,
		Status:      "paid",
	})

	return &shop.BuyWithBonusResponse{
		Success: true,
		Message: "Всё сработало!",
	}, nil
}

func (s *ShopServer) CreateUser(ctx context.Context, req *shop.CreateUserRequest) (*shop.CreateUserResponse, error) {
	log.Printf("Пришёл запрос на создание пользователя: Name=%s", req.Name)
	userId, err := s.Repo.CreateUser(req.Name)
	if err != nil {
		log.Println("Ошибка при создании пользователя", err)
		return &shop.CreateUserResponse{Success: false, Message: "Ошибка при создании пользователя"}, nil
	}
	_, err = s.Repo.CreateBank(userId)
	if err != nil {
		log.Println("Не удалось создать аккаунт", err)
		return &shop.CreateUserResponse{Success: false, Message: "Не удалось создать аккаунт"}, nil
	}
	return &shop.CreateUserResponse{
		Success: true,
		Message: "Пользователь успешно создан!",
		UserId:  int32(userId),
	}, nil
}

func (s *ShopServer) CreateCard(ctx context.Context, req *shop.CreateCardRequest) (*shop.CreateCardResponse, error) {
	log.Printf("Запрос на создание карты: UserID=%d", req.UserId)

	var bankID int
	var err error
	userId := int(req.UserId)

	if userId == 0 {
		bankID, err = s.Repo.CreateBank(0)
		if err != nil {
			return &shop.CreateCardResponse{Success: false, Message: "Не удалось создать счет"}, nil
		}
	} else {
		bankID, err = s.Repo.GetBankByUserId(userId)
		if err != nil {
			return &shop.CreateCardResponse{Success: false, Message: "Счет пользователя не найден"}, nil
		}
	}
	cardNumber := fmt.Sprintf("%016d", rand.Int63n(1e16))
	cardID, err := s.Repo.CreateCard(bankID, cardNumber)
	if err != nil {
		return &shop.CreateCardResponse{Success: false, Message: "Ошибка при привязке карты"}, nil
	}
	return &shop.CreateCardResponse{
		Success: true,
		Message: "Карта успешно создана!",
		CardId:  int32(cardID),
		Number:  cardNumber,
	}, nil
}

func (s *ShopServer) TopUp(ctx context.Context, req *shop.TopUpRequest) (*shop.TopUpResponse, error) {
	log.Printf("Пришёл запрос на пополнение карты: CardId=%d", req.CardId)
	card, err := s.Repo.GetCard(int(req.CardId))
	if err != nil {
		return &shop.TopUpResponse{Success: false, Message: "Карта не найдена"}, nil
	}
	if req.Valute == "RUB" {
		err = s.Repo.UpdateBalance(card.BankID, int(req.Amount), "RUB")
		if err != nil {
			return &shop.TopUpResponse{Success: false, Message: "Не удалось обновить баланс"}, nil
		}
	} else if req.Valute == "USD" {
		amount := int(req.Amount / 90)
		err = s.Repo.UpdateBalance(card.BankID, amount, "USD")
		if err != nil {
			return &shop.TopUpResponse{Success: false, Message: "Не удалось обновить баланс"}, nil
		}
	}
	return &shop.TopUpResponse{
		Success:    true,
		Message:    "Счет успешно пополнен!",
		NewBalance: req.Amount,
	}, nil
}

func (s *ShopServer) GetOrder(ctx context.Context, req *shop.GetOrderRequest) (*shop.GetOrderResponse, error) {
	log.Printf("Пришёл запрос на получение заказа: UserID=%d", req.UserId)
	order, err := s.Repo.GetOrder(int(req.UserId))
	if err != nil {
		return &shop.GetOrderResponse{Success: false, Message: "Заказ не найден"}, nil
	}
	return &shop.GetOrderResponse{
		Success: true,
		Message: "Заказ успешно получен!",
		Order:   fmt.Sprintf("ID: %d, UserID: %d, ProductID: %d, Amount: %d, Price: %d, UsedBonuses: %d, Status: %s", order.ID, order.UserID, order.ProductID, order.Amount, order.Price, order.UsedBonuses, order.Status),
	}, nil
}
func (s *ShopServer) LinkCard(ctx context.Context, req *shop.LinkCardRequest) (*shop.LinkCardResponse, error) {
	log.Printf("Пришёл запрос на активацию карты: UserID=%d", req.UserId)
	card, err := s.Repo.GetCard(int(req.CardId))
	if err != nil {
		return &shop.LinkCardResponse{Success: false, Message: "Не удалось найти карту"}, nil
	}
	oldBank, err := s.Repo.GetBankById(card.BankID)
	if err != nil {
		return &shop.LinkCardResponse{Success: false, Message: "Не удалось найти счет карты"}, nil
	}

	if oldBank.UserID != 0 {
		return &shop.LinkCardResponse{Success: false, Message: "Карта уже привязана к другому пользователю!"}, nil
	}
	userBankID, err := s.Repo.GetBankByUserId(int(req.UserId))
	if err != nil {
		return &shop.LinkCardResponse{Success: false, Message: "Не удалось найти счет пользователя"}, nil
	}
	err = s.Repo.UpdateCardBank(int(req.CardId), userBankID)
	if err != nil {
		return &shop.LinkCardResponse{Success: false, Message: "Не удалось перепривязать карту"}, nil
	}
	if oldBank.BalanceRub > 0 {
		s.Repo.UpdateBalance(userBankID, oldBank.BalanceRub, "RUB")
		s.Repo.UpdateBalance(oldBank.ID, -oldBank.BalanceRub, "RUB")
	}
	if oldBank.BalanceUsd > 0 {
		s.Repo.UpdateBalance(userBankID, oldBank.BalanceUsd, "USD")
		s.Repo.UpdateBalance(oldBank.ID, -oldBank.BalanceUsd, "USD")
	}
	return &shop.LinkCardResponse{
		Success: true,
		Message: "Карта успешно активирована, баланс перенесен!",
	}, nil
}
