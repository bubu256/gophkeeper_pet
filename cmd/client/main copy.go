// package main

// import (
// 	"context"
// 	"log"
// 	"strings"

// 	"github.com/bubu256/gophkeeper_pet/config"
// 	"github.com/bubu256/gophkeeper_pet/internal/proto/pb" // Путь к сгенерированному файлу протокола
// 	"github.com/joho/godotenv"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/metadata"
// )

// type Client struct {
// 	token  string
// 	client pb.GophKeeperServiceClient
// 	ctx    context.Context
// }

// func main() {

// 	godotenv.Load()

// 	cfg, err := config.GetClientConfig()
// 	if err != nil {
// 		log.Fatalf("configuration loading failed %v", err)
// 	}

// 	srvAddress := strings.Join([]string{cfg.ServerAddress, cfg.Port}, ":")
// 	conn, err := grpc.Dial(srvAddress, grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("Failed to dial server: %v", err)
// 	}
// 	defer conn.Close()

// 	// Создание клиента
// 	client := pb.NewGophKeeperServiceClient(conn)

// 	// Создание объекта конфигурации клиента
// 	// clientConfig := config.ClientConfig{
// 	// 	ServerAddress: "localhost:50051", // Адрес сервера gRPC
// 	// }

// 	// Установка соединения с сервером

// 	// Вызов методов сервера
// 	// registrationResponse, err := registerUser(client)
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to register user: %v", err)
// 	// }
// 	// log.Printf("Registration Response: %v", registrationResponse)

// 	authenticationResponse, err := authenticateUser(client)
// 	if err != nil {
// 		log.Fatalf("Failed to authenticate user: %v", err)
// 	}
// 	log.Printf("Authentication Response: %+v", authenticationResponse)
// 	token := authenticationResponse.Token

// 	// authorizationResponse, err := authorizeUser(client, authenticationResponse.Token)
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to authorize user: %v", err)
// 	// }
// 	// log.Printf("Authorization Response: %+v", authorizationResponse.Success)

// 	addDataResponse, err := addData(client, token)
// 	if err != nil {
// 		log.Fatalf("Failed to add data: %v", err)
// 	}
// 	log.Printf("AddData Response: %+v", addDataResponse)

// 	retrieveDataResponse, err := retrieveData(client, token)
// 	if err != nil {
// 		log.Printf("Failed to retrieve data: %v", err)
// 	}
// 	log.Printf("RetrieveData Response: %v", retrieveDataResponse)

// 	getInformationResponse, err := getInformation(client, token)
// 	if err != nil {
// 		log.Printf("Failed to get information: %v", err)
// 	}
// 	log.Printf("GetInformation Response: %v", getInformationResponse)
// }

// func registerUser(client pb.GophKeeperServiceClient) (*pb.RegistrationResponse, error) {
// 	request := &pb.RegistrationRequest{
// 		Username: "test_user",
// 		Password: "234",
// 	}

// 	response, err := client.Register(context.Background(), request)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return response, nil
// }

// func authenticateUser(client pb.GophKeeperServiceClient) (*pb.AuthenticationResponse, error) {
// 	request := &pb.AuthenticationRequest{
// 		Username: "test_user",
// 		Password: "234",
// 	}

// 	response, err := client.Authenticate(context.Background(), request)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return response, nil
// }

// func authorizeUser(client pb.GophKeeperServiceClient, token string) (*pb.AuthorizationResponse, error) {
// 	request := &pb.AuthorizationRequest{
// 		Token: token,
// 	}

// 	response, err := client.Authorize(SetTokenContext(context.Background(), token), request)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return response, nil
// }

// func addData(client pb.GophKeeperServiceClient, token string) (*pb.AddDataResponse, error) {
// 	request := &pb.AddDataRequest{
// 		Data: &pb.MemoryCell{
// 			Info: &pb.InfoCell{
// 				Id:          1320947,
// 				DataType:    "sample-data-type",
// 				DataSize:    1024,
// 				Description: "Sample data description",
// 				OwnerId:     1,
// 			},
// 			Encrypted:     false,
// 			KeyValuePairs: map[string]string{"login": "value1", "password": "value2"},
// 			BinaryData:    []byte("sample-binary-data"),
// 			FileName:      "sample-file-name",
// 		},
// 	}

// 	response, err := client.AddData(SetTokenContext(context.Background(), token), request)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return response, nil
// }

// func retrieveData(client pb.GophKeeperServiceClient, token string) (*pb.RetrieveDataResponse, error) {
// 	request := &pb.RetrieveDataRequest{
// 		Ids: []int64{140},
// 	}

// 	response, err := client.RetrieveData(SetTokenContext(context.Background(), token), request)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return response, nil
// }

// func getInformation(client pb.GophKeeperServiceClient, token string) (*pb.GetInformationResponse, error) {
// 	request := &pb.GetInformationRequest{}

// 	response, err := client.GetInformation(SetTokenContext(context.Background(), token), request)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return response, nil
// }

// func SetTokenContext(ctx context.Context, token string) context.Context {
// 	md := metadata.New(map[string]string{"token": token})
// 	return metadata.NewOutgoingContext(ctx, md)
// }
