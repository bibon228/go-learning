package main

import (
	shop "db/proto"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	InitDB()
	go autoUpdateOrders()
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer()
		worker := &PostgresRepo{db: db}
		shop.RegisterShopServiceServer(s, &ShopServer{Repo: worker})
		reflection.Register(s)

		fmt.Println("gRPC сервер запущен на :50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	http.HandleFunc("/create_user", handleCreateUser)
	http.HandleFunc("/create_card", handleCreateCard)
	http.HandleFunc("/top_up", handleTopUp)
	http.HandleFunc("/link_card", handleLinkCard)
	http.HandleFunc("/create_product", handleCreateProduct)
	http.HandleFunc("/buy_product", handleBuyProduct)
	http.HandleFunc("/buy_with_bonus", handleBuyWithBonus)

	fmt.Println("HTTP сервер запущен на :8080")
	http.ListenAndServe(":8080", nil)

}
