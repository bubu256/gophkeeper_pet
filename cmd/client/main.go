package main

import (
	"context"
	"log"

	"github.com/bubu256/gophkeeper_pet/config"
	"github.com/bubu256/gophkeeper_pet/internal/proto/pb" // Путь к сгенерированному файлу протокола
	"google.golang.org/grpc"
)

func main() {
	// Создание объекта конфигурации клиента
	clientConfig := config.ClientConfig{
		ServerAddress: "localhost:50051", // Адрес сервера gRPC
	}

	// Установка соединения с сервером
	conn, err := grpc.Dial(clientConfig.ServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// Создание клиента
	client := pb.NewGophKeeperServiceClient(conn)

	// Вызов методов сервера
	registrationResponse, err := registerUser(client)
	if err != nil {
		log.Fatalf("Failed to register user: %v", err)
	}
	log.Printf("Registration Response: %v", registrationResponse)

	authenticationResponse, err := authenticateUser(client)
	if err != nil {
		log.Fatalf("Failed to authenticate user: %v", err)
	}
	log.Printf("Authentication Response: %v", authenticationResponse)

	authorizationResponse, err := authorizeUser(client)
	if err != nil {
		log.Fatalf("Failed to authorize user: %v", err)
	}
	log.Printf("Authorization Response: %v", authorizationResponse)

	addDataResponse, err := addData(client)
	if err != nil {
		log.Fatalf("Failed to add data: %v", err)
	}
	log.Printf("AddData Response: %v", addDataResponse)

	retrieveDataResponse, err := retrieveData(client)
	if err != nil {
		log.Fatalf("Failed to retrieve data: %v", err)
	}
	log.Printf("RetrieveData Response: %v", retrieveDataResponse)

	getInformationResponse, err := getInformation(client)
	if err != nil {
		log.Fatalf("Failed to get information: %v", err)
	}
	log.Printf("GetInformation Response: %v", getInformationResponse)
}

func registerUser(client pb.GophKeeperServiceClient) (*pb.RegistrationResponse, error) {
	request := &pb.RegistrationRequest{
		Username: "example_user",
		Password: "password123",
	}

	response, err := client.Register(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func authenticateUser(client pb.GophKeeperServiceClient) (*pb.AuthenticationResponse, error) {
	request := &pb.AuthenticationRequest{
		Username: "example_user",
		Password: "password123",
	}

	response, err := client.Authenticate(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func authorizeUser(client pb.GophKeeperServiceClient) (*pb.AuthorizationResponse, error) {
	request := &pb.AuthorizationRequest{
		Token: "sample-token",
	}

	response, err := client.Authorize(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func addData(client pb.GophKeeperServiceClient) (*pb.AddDataResponse, error) {
	request := &pb.AddDataRequest{
		Data: []*pb.MemoryCell{
			{
				Info: &pb.InfoCell{
					Id:          "data-id-1",
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

	response, err := client.AddData(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func retrieveData(client pb.GophKeeperServiceClient) (*pb.RetrieveDataResponse, error) {
	request := &pb.RetrieveDataRequest{
		Id: []string{"data-id-1"},
	}

	response, err := client.RetrieveData(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func getInformation(client pb.GophKeeperServiceClient) (*pb.GetInformationResponse, error) {
	request := &pb.GetInformationRequest{}

	response, err := client.GetInformation(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
