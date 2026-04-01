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
	worker := &PostgresRepo{db: db}
	httpServer := &HttpHandler{Repo: worker}

	go httpServer.autoUpdateOrders()
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()

		shop.RegisterShopServiceServer(s, &ShopServer{Repo: worker})
		reflection.Register(s)
		fmt.Println("gRPC сервер запущен на :50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	http.HandleFunc("/create_user", httpServer.handleCreateUser)
	http.HandleFunc("/create_card", httpServer.handleCreateCard)
	http.HandleFunc("/top_up", httpServer.handleTopUp)
	http.HandleFunc("/link_card", httpServer.handleLinkCard)
	http.HandleFunc("/create_product", httpServer.handleCreateProduct)
	http.HandleFunc("/buy_product", httpServer.handleBuyProduct)
	http.HandleFunc("/buy_with_bonus", httpServer.handleBuyWithBonus)

	fmt.Println("HTTP сервер запущен на :8080")
	http.ListenAndServe(":8080", nil)

}
