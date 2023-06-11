package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/bubu256/gophkeeper_pet/config"
	"github.com/bubu256/gophkeeper_pet/internal/proto/ghandlers"
	"github.com/bubu256/gophkeeper_pet/pkg/keeper"

	"github.com/bubu256/gophkeeper_pet/internal/goph"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg, err := config.GetServerConfig()
	if err != nil {
		log.Fatalf("configuration loading failed %v", err)
	}
	storage, err := keeper.New(cfg)
	if err != nil {
		log.Fatalf("Storage creation failed %v", err)
	}
	logic := goph.New(storage, cfg)

	server := ghandlers.New(logic, cfg)

	// Запуск сервера
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}
	log.Printf("Server is listening on port %s", cfg.Port)

	// Запуск grpc сервера в отдельной goroutine
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Ожидание сигнала прерывания (Ctrl+C)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	// Остановка сервера
	server.Stop()
	log.Println("Server stopped")
}
