package main

import (
	"context"
	"log"
	"time"

	"leads/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Подключаемся к серверу
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer conn.Close()

	client := pb.NewLeadServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// 1. Тестируем создание лида
	log.Println("--- 1. Создание Лида ---")
	createResp, err := client.CreateLead(ctx, &pb.CreateLeadRequest{
		Name:   "Иван Иванов",
		Phone:  "+79998887766",
		Email:  "ivan@example.com",
		Source: "Website",
	})
	if err != nil {
		log.Fatalf("Ошибка при создании лида: %v", err)
	}
	log.Printf("Лид успешно создан. ID: %s", createResp.GetId())

	createdLeadId := createResp.GetId()

	// 2. Тестируем получение лида по ID
	log.Println("\n--- 2. Получение Лида по ID ---")
	getResp, err := client.GetLead(ctx, &pb.GetLeadRequest{
		Id: createdLeadId,
	})
	if err != nil {
		log.Printf("Ошибка при получении лида: %v", err)
	} else {
		log.Printf("Получен лид: %+v", getResp.GetLead())
	}

	// 3. Тестируем обновление статуса
	log.Println("\n--- 3. Обновление статуса Лида ---")
	updateResp, err := client.UpdateLeadStatus(ctx, &pb.UpdateLeadStatusRequest{
		Id:     createdLeadId,
		Status: "В работе",
	})
	if err != nil {
		log.Printf("Ошибка при обновлении статуса: %v", err)
	} else {
		log.Printf("Статус обновлен? Успех=%v", updateResp)
	}

	// 4. Тестируем список лидов
	log.Println("\n--- 4. Список всех Лидов ---")
	listResp, err := client.ListLeads(ctx, &pb.ListLeadsRequest{
		Limit: 10,
	})
	if err != nil {
		log.Fatalf("Ошибка при получении списка: %v", err)
	}
	for i, lead := range listResp.GetLeads() {
		log.Printf("[%d] ID: %s, Name: %s, Status: %s", i+1, lead.GetId(), lead.GetName(), lead.GetStatus())
	}
}
