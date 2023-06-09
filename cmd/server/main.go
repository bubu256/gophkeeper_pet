package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/bubu256/gophkeeper_pet/config"
	"github.com/bubu256/gophkeeper_pet/internal/proto/ghandlers"
)

func main() {
	// Создание объекта реализующего бизнес-логику приложения
	// businessLogic := NewBusinessLogic() // Замените на вашу реализацию

	// Создание серверной конфигурации gRPC
	serverConfig := config.ServerConfig{
		// Добавьте необходимые опции сервера
	}

	// Создание нового gRPC сервера с использованием объекта реализации бизнес-логики
	server := ghandlers.New(serverConfig)

	// Запуск сервера на заданном порту
	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", port, err)
	}
	log.Printf("Server is listening on port %d", port)

	// Запуск сервера в отдельной goroutine
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
