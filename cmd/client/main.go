package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bubu256/gophkeeper_pet/config"
	"github.com/bubu256/gophkeeper_pet/internal/proto/pb"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	token  string
	client pb.GophKeeperServiceClient
	ctx    context.Context
}

func main() {
	godotenv.Load()

	cfg, err := config.GetClientConfig()
	if err != nil {
		log.Fatalf("configuration loading failed: %v", err)
	}

	srvAddress := strings.Join([]string{cfg.ServerAddress, cfg.Port}, ":")
	conn, err := grpc.Dial(srvAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// Создание клиента
	client := pb.NewGophKeeperServiceClient(conn)

	cli := &Client{
		client: client,
		ctx:    context.Background(),
	}

	// Обработка сигналов для возможности выхода по запросу пользователя
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("Выход из приложения")
		os.Exit(0)
	}()

	// Запуск интерактивного меню
	cli.runMenu()
}

func (c *Client) runMenu() {
	for {
		fmt.Println("Меню:  1. Регистрация;  2. Авторизация;  3. Получение информации;  4. Получение данных по InfoID;       0. Выход")

		var choice string
		fmt.Print("Выберите пункт меню: ")
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			c.register()
		case "2":
			c.authenticate()
		case "3":
			c.retrieveInformation()
		case "4":
			c.retrieveDataByID()
		case "0":
			fmt.Println("Выход из приложения")
			os.Exit(0)
		default:
			fmt.Println("Некорректный выбор")
		}

		fmt.Println()
	}
}

func (c *Client) register() {
	// запрос данных у пользователя
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите логин: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Ошибка при чтении логина: %v", err)
	}
	username = strings.TrimSpace(username)

	fmt.Print("Введите пароль: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Ошибка при чтении пароля: %v", err)
	}
	password = strings.TrimSpace(password)

	// создание запроса регистрации
	request := &pb.RegistrationRequest{
		Username: username,
		Password: password,
	}

	// отправка запроса на сервер
	response, err := c.client.Register(c.ctx, request)
	if err != nil {
		log.Fatalf("Ошибка при отправке запроса регистрации: %v", err)
	}

	// вывод информации о результате
	if response.Success {
		fmt.Println("Регистрация прошла успешно.")
	} else {
		fmt.Println("Ошибка при регистрации.")
	}
}

func (c *Client) authenticate() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите логин: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения ввода пользователя:", err)
		return
	}
	username = strings.TrimSpace(username)

	fmt.Print("Введите пароль: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения ввода пользователя:", err)
		return
	}
	password = strings.TrimSpace(password)

	request := &pb.AuthenticationRequest{
		Username: username,
		Password: password,
	}

	response, err := c.client.Authenticate(c.ctx, request)
	if err != nil {
		fmt.Println("Ошибка при аутентификации:", err)
		return
	}

	c.token = response.Token
	c.ctx = SetTokenContext(c.ctx, c.token)

	fmt.Println("Аутентификация прошла успешно.")
}

func (c *Client) retrieveInformation() {
	if c.token == "" {
		fmt.Println("Необходима авторизация")
		return
	}

	// Логика получения информации о сохраненных данных
	// Используйте методы клиента (c.client) и контекст (c.ctx) для вызова gRPC-запросов
	// Обработка ошибок и вывод результатов
}

func (c *Client) retrieveDataByID() {
	if c.token == "" {
		fmt.Println("Необходима авторизация")
		return
	}

	// Логика получения данных по введенному InfoID
	// Используйте методы клиента (c.client) и контекст (c.ctx) для вызова gRPC-запросов
	// Обработка ошибок и вывод результатов
}

func SetTokenContext(ctx context.Context, token string) context.Context {
	md := metadata.New(map[string]string{"token": token})
	return metadata.NewOutgoingContext(ctx, md)
}
