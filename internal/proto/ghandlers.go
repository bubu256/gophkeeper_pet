// Package ghandler - реализует работу grpc сервиса и содержит имплементации интерфейсов описанных proto файле
package ghandler

import (
	"context"

	pb "github.com/bubu256/gophkeeper_pet/internal/proto/pb" // Путь к сгенерированному файлу протокола

	"google.golang.org/grpc"
)

// HandlerService представляет собой структуру, реализующую интерфейсы сервера gRPC.
type HandlerService struct {
	appLogic AppLogic // Объект, реализующий бизнес-логику приложения
}

// New создает новый объект HandlerService с указанным объектом AppLogic и возвращает ссылку на grpc.Server.
func New(appLogic AppLogic, config ServerConfig) (*grpc.Server, error) {
	// Создание gRPC сервера
	server := grpc.NewServer()

	// Инициализация HandlerService
	handler := &HandlerService{
		appLogic: appLogic,
	}

	// Регистрация HandlerService в gRPC сервере
	pb.RegisterGophKeeperServiceServer(server, handler)

	return server, nil
}

// Register реализует метод Register интерфейса GophKeeperServiceServer.
func (h *HandlerService) Register(ctx context.Context, req *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
	// Реализация логики регистрации пользователя
	// Используйте h.appLogic для вызова соответствующего метода бизнес-логики
	// ...

	return &pb.RegistrationResponse{
		Success: true, // Пример возвращаемого результата
	}, nil
}

// Authenticate реализует метод Authenticate интерфейса GophKeeperServiceServer.
func (h *HandlerService) Authenticate(ctx context.Context, req *pb.AuthenticationRequest) (*pb.AuthenticationResponse, error) {
	// Реализация логики аутентификации пользователя
	// Используйте h.appLogic для вызова соответствующего метода бизнес-логики
	// ...

	return &pb.AuthenticationResponse{
		Token: "example_token", // Пример возвращаемого результата
	}, nil
}

// Authorize реализует метод Authorize интерфейса GophKeeperServiceServer.
func (h *HandlerService) Authorize(ctx context.Context, req *pb.AuthorizationRequest) (*pb.AuthorizationResponse, error) {
	// Реализация логики авторизации пользователя
	// Используйте h.appLogic для вызова соответствующего метода бизнес-логики
	// ...

	return &pb.AuthorizationResponse{
		Success: true, // Пример возвращаемого результата
	}, nil
}

// AddData реализует метод AddData интерфейса GophKeeperServiceServer.
func (h *HandlerService) AddData(ctx context.Context, req *pb.AddDataRequest) (*pb.AddDataResponse, error) {
	// Реализация логики добавления данных
	// Используйте h.appLogic для вызова соответствующего метода бизнес-логики
	// ...

	return &pb.AddDataResponse{
		Id: []string{"data_id1", "data_id2"}, // Пример возвращаемого результата
	}, nil
}

// RetrieveData реализует метод RetrieveData интерфейса GophKeeperServiceServer.
func (h *HandlerService) RetrieveData(ctx context.Context, req *pb.RetrieveDataRequest) (*pb.RetrieveDataResponse, error) {
	// Реализация логики запроса данных
	// Используйте h.appLogic для вызова соответствующего метода бизнес-логики
	// ...

	return &pb.RetrieveDataResponse{
		Data: []*pb.MemoryCell{ // Пример возвращаемого результата
			{
				Info: &pb.InfoCell{
					Id:          "data_id",
					DataType:    "data_type",
					DataSize:    123,
					Description: "data_description",
					OwnerId:     "owner_id",
				},
				Encrypted:     false,
				KeyValuePairs: map[string]string{"key": "value"},
				BinaryData:    []byte("binary_data"),
				FileName:      "file_name",
			},
		},
	}, nil
}

// GetInformation реализует метод GetInformation интерфейса GophKeeperServiceServer.
func (h *HandlerService) GetInformation(ctx context.Context, req *pb.GetInformationRequest) (*pb.GetInformationResponse, error) {
	// Реализация логики запроса информации о данных
	// Используйте h.appLogic для вызова соответствующего метода бизнес-логики
	// ...

	return &pb.GetInformationResponse{
		Info: []*pb.InfoCell{ // Пример возвращаемого результата
			{
				Id:          "data_id",
				DataType:    "data_type",
				DataSize:    123,
				Description: "data_description",
				OwnerId:     "owner_id",
			},
		},
	}, nil
}
