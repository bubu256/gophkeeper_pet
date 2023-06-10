// Package ghandler - реализует работу grpc сервиса и содержит имплементации интерфейсов описанных proto файле
package ghandlers

import (
	"context"
	"log"

	"github.com/bubu256/gophkeeper_pet/config"
	pb "github.com/bubu256/gophkeeper_pet/internal/proto/pb" // Путь к сгенерированному файлу протокола

	"google.golang.org/grpc"
)

// Goph - интерфейс бизнес логики приложения
type Goph interface {
	RegisterUser(ctx context.Context, request *pb.RegistrationRequest) (*pb.RegistrationResponse, error)
	AuthenticateUser(ctx context.Context, request *pb.AuthenticationRequest) (*pb.AuthenticationResponse, error)
	AuthorizeUser(ctx context.Context, request *pb.AuthorizationRequest) (*pb.AuthorizationResponse, error)
	AddData(ctx context.Context, request *pb.AddDataRequest) (*pb.AddDataResponse, error)
	RetrieveData(ctx context.Context, request *pb.RetrieveDataRequest) (*pb.RetrieveDataResponse, error)
	GetInformation(ctx context.Context, request *pb.GetInformationRequest) (*pb.GetInformationResponse, error)
}

// HandlerService представляет собой структуру, реализующую интерфейсы сервера gRPC.
type HandlerService struct {
	pb.UnimplementedGophKeeperServiceServer
	businessLogic Goph
	cfg           config.ServerConfig
}

// New создает новый объект HandlerService и возвращает ссылку на grpc.Server.
func New(serverConfig config.ServerConfig) *grpc.Server {
	handler := &HandlerService{
		// businessLogic: businessLogic,
		cfg: serverConfig,
	}

	server := grpc.NewServer()
	pb.RegisterGophKeeperServiceServer(server, handler)

	return server
}

// Register реализует метод Register интерфейса GophKeeperServiceServer.
func (h *HandlerService) Register(ctx context.Context, request *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
	log.Println("Received Register request")
	log.Printf("Username: %s, Password: %s", request.Username, request.Password)

	response := &pb.RegistrationResponse{
		Success: true,
	}

	return response, nil
}

// Authenticate реализует метод Authenticate интерфейса GophKeeperServiceServer.
func (h *HandlerService) Authenticate(ctx context.Context, request *pb.AuthenticationRequest) (*pb.AuthenticationResponse, error) {
	log.Println("Received Authenticate request")
	log.Printf("Username: %s, Password: %s", request.Username, request.Password)

	response := &pb.AuthenticationResponse{
		Token: "sample-token",
	}

	return response, nil
}

// Authorize реализует метод Authorize интерфейса GophKeeperServiceServer.
func (h *HandlerService) Authorize(ctx context.Context, request *pb.AuthorizationRequest) (*pb.AuthorizationResponse, error) {
	log.Println("Received Authorize request")
	log.Printf("Token: %s", request.Token)

	response := &pb.AuthorizationResponse{
		Success: true,
	}

	return response, nil
}

// AddData реализует метод AddData интерфейса GophKeeperServiceServer.
func (h *HandlerService) AddData(ctx context.Context, request *pb.AddDataRequest) (*pb.AddDataResponse, error) {
	log.Println("Received AddData request")
	log.Printf("Data: %v", request.Data)

	response := &pb.AddDataResponse{
		Id: []int32{1, 2},
	}

	return response, nil
}

// RetrieveData реализует метод RetrieveData интерфейса GophKeeperServiceServer.
func (h *HandlerService) RetrieveData(ctx context.Context, request *pb.RetrieveDataRequest) (*pb.RetrieveDataResponse, error) {
	log.Println("Received RetrieveData request")
	log.Printf("Ids: %v", request.Id)

	response := &pb.RetrieveDataResponse{
		Data: []*pb.MemoryCell{
			{
				Info: &pb.InfoCell{
					Id:          1,
					DataType:    "sample-data-type",
					DataSize:    1024,
					Description: "Sample data description",
					OwnerId:     "user-id",
				},
				Encrypted:     false,
				KeyValuePairs: map[string]string{"key1": "value1", "key2": "value2"},
				BinaryData:    []byte("sample-binary-data"),
				FileName:      "sample-file-name",
			},
		},
	}

	return response, nil
}

// GetInformation реализует метод GetInformation интерфейса GophKeeperServiceServer.
func (h *HandlerService) GetInformation(ctx context.Context, request *pb.GetInformationRequest) (*pb.GetInformationResponse, error) {
	log.Println("Received GetInformation request")

	response := &pb.GetInformationResponse{
		Info: []*pb.InfoCell{
			{
				Id:          1,
				DataType:    "sample-data-type",
				DataSize:    1024,
				Description: "Sample data description",
				OwnerId:     "user-id",
			},
		},
	}

	return response, nil
}
