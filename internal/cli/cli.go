// Package cli - организует консольный интерфейс для клиентской части приложения. Отправляет запросы на сервер или использует локальный кэш.
package cli

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bubu256/gophkeeper_pet/internal/proto/pb"
	"github.com/bubu256/gophkeeper_pet/pkg/keeper/client/storage"
	"google.golang.org/grpc/metadata"
)

// Cli - структура для реализации работы меню приложения и его взаимодействия с grpc серверомю
type Cli struct {
	token   string
	client  pb.GophKeeperServiceClient
	ctx     context.Context
	storage *storage.Storage
	Cancel  context.CancelFunc
}

// NewCli - возвращает экземпляр Cli
func NewCli(client pb.GophKeeperServiceClient, ctx context.Context, storage *storage.Storage) *Cli {
	ctx, cancel := context.WithCancel(ctx)
	return &Cli{
		client:  client,
		ctx:     ctx,
		storage: storage,
		Cancel:  cancel,
	}
}

// RunMenu - запускает интерактивное меню приложения
func (c *Cli) RunMenu() {
	for {
		fmt.Println("Меню:  1. Регистрация;  2. Авторизация;  3. Получение информации;  4. Получение данных по InfoID;  5. Сохранить новые данные;     0. Выход")

		var choice string
		fmt.Print("Выберите пункт меню: ")
		fmt.Scanln(&choice)

		select {
		case <-c.ctx.Done():
			fmt.Println("Выход из приложения")
			return
		default:
			switch choice {
			case "1":
				c.Register()
			case "2":
				c.Authenticate()
			case "3":
				c.RetrieveInformation()
			case "4":
				c.RetrieveDataByID()
			case "5":
				c.AddData()
			case "0":
				c.Dump()
				fmt.Println("Выход из приложения")
				return
			default:
				fmt.Println("Некорректный выбор")
			}

			fmt.Println()
		}
	}
}

// Register - регистрация пользователя
func (c *Cli) Register() {
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
		fmt.Println("- Регистрация прошла успешно.")
	} else {
		fmt.Println("- Ошибка при регистрации.")
	}
}

// Authenticate - аутентификация пользователя
func (c *Cli) Authenticate() {
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
		fmt.Println("- Ошибка при аутентификации:", err)
		return
	}

	// сохраняем токен и пишем в контекст
	c.token = response.Token
	c.ctx = SetTokenContext(c.ctx, c.token)
	fmt.Println("- Аутентификация прошла успешно.")
	// синхронизация клиента
	c.Sync()
}

// RetrieveInformation - получение мета информации о данных пользователя
func (c *Cli) RetrieveInformation() {
	response, err := c.client.GetInformation(c.ctx, &pb.GetInformationRequest{})
	if err != nil {
		fmt.Println("Ошибка при получении информации:", err)
		return
	}

	fmt.Println("\tСписок сохраненных данных:")
	fmt.Println("\t-------------------------------")
	for _, info := range response.Info {
		fmt.Printf("\tID: %d\n", info.Id)
		fmt.Printf("\tТип данных: %s\n", info.DataType)
		fmt.Printf("\tРазмер данных: %d\n", info.DataSize)
		fmt.Printf("\tОписание: %s\n", info.Description)
		fmt.Printf("\tID владельца: %d\n", info.OwnerId)
		fmt.Println("\t-------------------------------")
	}
}

// RetrieveDataByID - получение данных по InfoID из мета информации.
// Мету информацию пользователь может получить с помощью метода RetrieveInformation
func (c *Cli) RetrieveDataByID() {
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
		fmt.Println("- Данные не найдены")
		return
	}

	// Вывод данных MemoryCell
	data := response.Data[0]
	fmt.Println("\tПолученные данные:")
	// fmt.Printf("ID: %d\n", data.Id)
	fmt.Printf("\tInfoID: %d\n", data.Info.Id)
	fmt.Printf("\tТип данных: %s\n", data.Info.DataType)
	fmt.Printf("\tОписание: %s\n", data.Info.Description)
	fmt.Printf("\tЗашифровано: %v\n", data.Encrypted)
	fmt.Println("\tКлючи:")
	for key, value := range data.KeyValuePairs {
		fmt.Printf("  %s: %s\n", key, value)
	}
	fmt.Printf("\tBinaryData: %v байт\n", len(data.BinaryData))
	fmt.Printf("\tFileName: %s\n", data.FileName)
	if (len(data.BinaryData) > 0) && (data.FileName != "") {
		fmt.Print("\tДля сохранения файла укажите путь к папке: ")
		var pathdir string
		fmt.Scanln(&pathdir)
		if pathdir == "" {
			fmt.Print("\tФайл не сохранен.")
		} else {
			err := WriteBytesToFile(data.BinaryData, pathdir, data.FileName)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Println("\tФайл сохранен.")
			}
		}
	}
}

// Добавление новых данных пользователя на сервер
func (c *Cli) AddData() {
	// Запрос данных у пользователя
	data := pb.MemoryCell{Info: &pb.InfoCell{}}

	fmt.Println("Введите данные для сохранения:")

	fmt.Print("Тип данных: ")
	fmt.Scanln(&data.Info.DataType)

	fmt.Print("Описание: ")
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
		var key string
		fmt.Scanln(&key)
		if key == "" {
			break
		}
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
		fmt.Println("- Ошибка при добавлении данных:", err)
		return
	}

	fmt.Println("- Данные успешно добавлены. ID данных:", response.Id)
}

// GetAllData - возвращает все данные пользователя
func (c *Cli) GetAllData() []*pb.MemoryCell {
	response, err := c.client.GetInformation(c.ctx, &pb.GetInformationRequest{})
	if err != nil {
		fmt.Println("Ошибка при получении информации:", err)
		return nil
	}
	// создаем слайс из всех InfoID пользователя
	infoIDs := make([]int64, len(response.Info))
	for i, info := range response.Info {
		infoIDs[i] = info.Id
	}
	// получаем все данные пользователя
	responseData, err := c.client.RetrieveData(c.ctx, &pb.RetrieveDataRequest{Ids: infoIDs})
	if err != nil {
		fmt.Println("Ошибка при получении данных:", err)
		return nil
	}
	// Проверка наличия данных
	if len(responseData.Data) == 0 {
		fmt.Println("- Данные не найдены")
		return nil
	}
	return responseData.Data
}

// Sync - синхронизует данные пользователя на клиенте и предлагает пользователю сохранить резервную копию
func (c *Cli) Sync() {
	fmt.Print("- Синхронизация... ")
	data := c.GetAllData()
	c.storage.Sync(data)
	fmt.Println("успешна.")
	c.Dump()
}

// Dump - предлагает сохранение резервной копии на диск.
// Потребуется ввод пароля для шифрования данных.
func (c *Cli) Dump() {
	reader := bufio.NewReader(os.Stdin)
	// Запрашиваем у пользователя требуется ли сохранение резервной копии
	fmt.Print("Хотите сохранить резервную копию? (y/n): ")
	saveBackup, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения ввода пользователя:", err)
		return
	}
	saveBackup = strings.TrimSpace(saveBackup)

	// Проверяем ответ пользователя
	if saveBackup == "y" || saveBackup == "yes" {
		// запрашиваем пароль для шифрования и сохранения дампа
		fmt.Print("Введите пароль: ")
		password, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения ввода пользователя:", err)
			return
		}
		password = strings.TrimSpace(password)

		// делаем дамп
		err = c.storage.DumpToFile(password)
		if err != nil {
			log.Printf("Failed to save data to disk: %v", err)
			return
		} else {
			fmt.Println("- Резервная копия сохранена.")
		}
	}
}

// SetTokenContext - устанавливает токен в контект
func SetTokenContext(ctx context.Context, token string) context.Context {
	md := metadata.New(map[string]string{"token": token})
	return metadata.NewOutgoingContext(ctx, md)
}

// ReadFileBytes - чтения файла в массив байт
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

// WriteBytesToFile - сохраняет массив байт в файл.
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

// GetFileNameFromPath - получает имя файла из строки с полным путем.
func GetFileNameFromPath(filePath string) (string, error) {
	fileName := filepath.Base(filePath)
	if fileName == "." || fileName == string(filepath.Separator) {
		return "", fmt.Errorf("некорректный путь к файлу: %s", filePath)
	}
	return fileName, nil
}
