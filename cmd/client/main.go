package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
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
		fmt.Println("Меню:  1. Регистрация;  2. Авторизация;  3. Получение информации;  4. Получение данных по InfoID;  5. Сохранить новые данные;     0. Выход")

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
		case "5":
			c.addData()
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
	response, err := c.client.GetInformation(c.ctx, &pb.GetInformationRequest{})
	if err != nil {
		fmt.Println("Ошибка при получении информации:", err)
		return
	}

	fmt.Println("Список сохраненных данных:")
	fmt.Println("-------------------------------")
	for _, info := range response.Info {
		fmt.Printf("ID: %d\n", info.Id)
		fmt.Printf("Тип данных: %s\n", info.DataType)
		fmt.Printf("Размер данных: %d\n", info.DataSize)
		fmt.Printf("Описание: %s\n", info.Description)
		fmt.Printf("ID владельца: %d\n", info.OwnerId)
		fmt.Println("-------------------------------")
	}
}

func (c *Client) retrieveDataByID() {
	// Запрос InfoID у пользователя
	fmt.Print("Введите InfoID для получения данных: ")
	var infoID int64
	_, err := fmt.Scanln(&infoID)
	if err != nil {
		fmt.Println("Ошибка при чтении InfoID:", err)
		return
	}

	request := &pb.RetrieveDataRequest{
		Ids: []int64{infoID},
	}

	response, err := c.client.RetrieveData(c.ctx, request)
	if err != nil {
		fmt.Println("Ошибка при получении данных:", err)
		return
	}

	// Проверка наличия данных
	if len(response.Data) == 0 {
		fmt.Println("Данные не найдены")
		return
	}

	// Вывод данных MemoryCell
	data := response.Data[0]
	fmt.Println("Полученные данные:")
	// fmt.Printf("ID: %d\n", data.Id)
	fmt.Printf("InfoID: %d\n", data.Info.Id)
	fmt.Printf("Тип данных: %s\n", data.Info.DataType)
	fmt.Printf("Описание: %s\n", data.Info.Description)
	fmt.Printf("Зашифровано: %v\n", data.Encrypted)
	fmt.Println("Ключи:")
	for key, value := range data.KeyValuePairs {
		fmt.Printf("  %s: %s\n", key, value)
	}
	fmt.Printf("BinaryData: %v байт\n", len(data.BinaryData))
	fmt.Printf("FileName: %s\n", data.FileName)
	if (len(data.BinaryData) > 0) && (data.FileName != "") {
		fmt.Print("Для сохранения файла укажите путь к папке: ")
		var pathdir string
		fmt.Scanln(&pathdir)
		if pathdir == "" {
			fmt.Print("Файл не сохранен.")
		} else {
			err := WriteBytesToFile(data.BinaryData, pathdir, data.FileName)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Println("Файл сохранен.")
			}
		}
	}
}

func (c *Client) addData() {
	// Запрос данных у пользователя
	data := pb.MemoryCell{Info: &pb.InfoCell{}}

	fmt.Println("Введите данные для сохранения:")

	fmt.Print("Тип данных: ")
	fmt.Scanln(&data.Info.DataType)

	fmt.Print("Описание: ")
	// fmt.Scanln(&data.Info.Description)
	// считываем всю введенную строку
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
	} else {
		data.Info.Description = strings.TrimSuffix(input, "\n")
	}

	fmt.Println("Вводите пары ключ и значение через пробел (пустую строку для завершения ввода):")
	data.KeyValuePairs = make(map[string]string)

	for {
		// fmt.Print("Ключ: ")
		var key string
		fmt.Scanln(&key)

		if key == "" {
			break
		}

		// fmt.Print("Значение: ")
		var value string
		fmt.Scanln(&value)

		data.KeyValuePairs[key] = value
	}

	var pathfile string
	fmt.Print("Файл: ")
	fmt.Scanln(&pathfile)
	if pathfile != "" {
		filename, err := GetFileNameFromPath(pathfile)
		if err != nil {
			log.Println(err)
		} else {
			data.FileName = filename
			filebytes, err := ReadFileBytes(pathfile)
			if err != nil {
				log.Println(err)
			} else {
				data.BinaryData = filebytes
			}
		}
	}

	request := &pb.AddDataRequest{
		Data: &data,
	}

	response, err := c.client.AddData(c.ctx, request)
	if err != nil {
		fmt.Println("Ошибка при добавлении данных:", err)
		return
	}

	fmt.Println("Данные успешно добавлены. ID данных:", response.Id)
}

func SetTokenContext(ctx context.Context, token string) context.Context {
	md := metadata.New(map[string]string{"token": token})
	return metadata.NewOutgoingContext(ctx, md)
}

func ReadFileBytes(filePath string) ([]byte, error) {
	// Проверяем существование файла
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("файл не существует: %s", filePath)
	}

	// Читаем содержимое файла
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %s", err)
	}

	return fileBytes, nil
}

func WriteBytesToFile(data []byte, directoryPath, fileName string) error {
	// Создаем папку, если она не существует
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		if err := os.MkdirAll(directoryPath, 0755); err != nil {
			return fmt.Errorf("ошибка создания папки: %s", err)
		}
	}

	// Создаем файл и записываем данные
	filePath := fmt.Sprintf("%s/%s", directoryPath, fileName)
	err := ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл: %s", err)
	}

	return nil
}

func GetFileNameFromPath(filePath string) (string, error) {
	fileName := filepath.Base(filePath)
	if fileName == "." || fileName == string(filepath.Separator) {
		return "", fmt.Errorf("некорректный путь к файлу: %s", filePath)
	}
	return fileName, nil
}
