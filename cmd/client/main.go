package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bubu256/gophkeeper_pet/config"
	"github.com/bubu256/gophkeeper_pet/internal/cli"
	"github.com/bubu256/gophkeeper_pet/internal/proto/pb"
	"github.com/bubu256/gophkeeper_pet/pkg/keeper/client/storage"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// Загрузка переменных из .env
	godotenv.Load()

	// создаем конфигурацию клиента
	cfg, err := config.GetClientConfig()
	if err != nil {
		log.Fatalf("configuration loading failed: %v", err)
	}

	// Создание grpc соединения
	srvAddress := strings.Join([]string{cfg.ServerAddress, cfg.Port}, ":")
	conn, err := grpc.Dial(srvAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// Создание клиента сервиса, локального хранилища и экземпляра cli приложения
	client := pb.NewGophKeeperServiceClient(conn)
	cliStorage := storage.NewStorage()
	cliclient := cli.NewCli(client, context.Background(), cliStorage)

	// Обработка сигналов для возможности выхода по запросу пользователя
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		// cliclient.Dump()
		fmt.Println("Выход из приложения")
		os.Exit(0)
	}()

	// Запуск интерактивного меню
	cliclient.RunMenu()
}
