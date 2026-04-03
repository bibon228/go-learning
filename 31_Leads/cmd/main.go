package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"

	grps "leads/internal/delivery/grpc"
	"leads/internal/repository"
	"leads/internal/service"
	"leads/pb"

	"google.golang.org/grpc"
)

func main() {
	// 1. Подключение к БД
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/crm?sslmode=disable")
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	// 2. Инициализация репозитория
	leadRepo := repository.NewPostgresLeadRepo(db)

	// 3. Инициализация сервиса
	leadService := service.NewLeadService(leadRepo)

	// 4. Инициализация gRPC сервера
	grpcServer := grpc.NewServer()
	leadHandler := grps.NewLeadHandler(leadService)

	// 5. Регистрация сервиса
	pb.RegisterLeadServiceServer(grpcServer, leadHandler)

	// 6. Запуск сервера
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}

	log.Println("Сервер запущен на порту 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Ошибка сервера:", err)
	}
}
