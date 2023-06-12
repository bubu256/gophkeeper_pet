// Package ghandler - реализует работу grpc сервиса и содержит имплементации интерфейсов описанных proto файле
package ghandlers

import (
	"context"
	"path/filepath"

	"github.com/bubu256/gophkeeper_pet/config"
	"github.com/bubu256/gophkeeper_pet/internal/goph"
	pb "github.com/bubu256/gophkeeper_pet/internal/proto/pb"
	"github.com/bubu256/gophkeeper_pet/internal/schema"
	"golang.org/x/exp/slices"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// HandlerService представляет собой структуру, реализующую интерфейсы сервера gRPC.
type HandlerService struct {
	pb.UnimplementedGophKeeperServiceServer
	gophKeeper goph.Goph
	cfg        config.ServerConfig
}

// New создает новый объект HandlerService и возвращает ссылку на grpc.Server.
func New(logic goph.Goph, serverConfig config.ServerConfig) *grpc.Server {
	handler := &HandlerService{
		gophKeeper: logic,
		cfg:        serverConfig,
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(handler.tokenInterceptor))
	pb.RegisterGophKeeperServiceServer(server, handler)

	return server
}

// Register реализует метод Register интерфейса GophKeeperServiceServer.
func (h *HandlerService) Register(ctx context.Context, request *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
	// Проверка, что пользователь с указанным именем не существует
	exists, err := h.gophKeeper.UserExists(request.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check user existence: %v", err)
	}
	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "user with username '%s' already exists", request.Username)
	}

	// Создание нового пользователя
	err = h.gophKeeper.CreateUser(request.Username, request.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	response := &pb.RegistrationResponse{
		Success: true,
	}

	return response, nil
}

// Authenticate реализует метод Authenticate интерфейса GophKeeperServiceServer.
func (h *HandlerService) Authenticate(ctx context.Context, request *pb.AuthenticationRequest) (*pb.AuthenticationResponse, error) {
	token, err := h.gophKeeper.Authenticate(request.Username, request.Password)
	// log.Println(token, err)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Authentication failed: %v", err)
	}

	response := &pb.AuthenticationResponse{
		Token: token,
	}

	return response, nil
}

func (h *HandlerService) Authorize(ctx context.Context, request *pb.AuthorizationRequest) (*pb.AuthorizationResponse, error) {
	valid, err := h.gophKeeper.CheckToken(request.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to authorize token: %v", err)
	}

	response := &pb.AuthorizationResponse{
		Success: valid,
	}

	return response, nil
}

// AddData реализует метод AddData интерфейса GophKeeperServiceServer.
func (h *HandlerService) AddData(ctx context.Context, request *pb.AddDataRequest) (*pb.AddDataResponse, error) {
	// log.Println("Received AddData request")
	// log.Printf("Data: %v", request.Data)
	token, ok := GetTokenFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "No token in context")
	}
	id, err := h.gophKeeper.GetUserIDFromToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unsuccessful getting id from token %v", err)
	}

	memoryCell := ConvertPBMemoryCellToSchema(request.Data)

	infoId, err := h.gophKeeper.SaveData(id, memoryCell)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save data %v", err)
	}
	response := &pb.AddDataResponse{
		Id: infoId,
	}
	return response, nil
}

// GetInformation реализует метод GetInformation интерфейса GophKeeperServiceServer.
func (h *HandlerService) GetInformation(ctx context.Context, request *pb.GetInformationRequest) (*pb.GetInformationResponse, error) {
	// log.Println("Received GetInformation request")
	// log.Printf("Info ID: %v", request.Id)
	token, ok := GetTokenFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "No token in context")
	}
	userID, err := h.gophKeeper.GetUserIDFromToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unsuccessful getting id from token %v", err)
	}

	info, err := h.gophKeeper.GetUserDataInfo(userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get information: %v", err)
	}

	pbInfo := make([]*pb.InfoCell, len(info))
	for i, d := range info {
		pbInfo[i] = ConvertSchemaInfoCellToPB(d)
	}

	response := &pb.GetInformationResponse{
		Info: pbInfo,
	}

	return response, nil
}

// RetrieveData реализует метод RetrieveData интерфейса GophKeeperServiceServer.
func (h *HandlerService) RetrieveData(ctx context.Context, request *pb.RetrieveDataRequest) (*pb.RetrieveDataResponse, error) {
	// log.Println("Received RetrieveData request")
	// log.Printf("Data ID: %v", request.Id)
	token, ok := GetTokenFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "No token in context")
	}
	usedID, err := h.gophKeeper.GetUserIDFromToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unsuccessful getting id from token %v", err)
	}

	infoIDs := request.Ids
	data, err := h.gophKeeper.GetUserMemoryData(usedID, infoIDs)
	// log.Println(len(data))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to retrieve data: %v", err)
	}

	pbData := make([]*pb.MemoryCell, len(data))
	for i, d := range data {
		pbData[i] = ConvertSchemaMemoryCellToPB(d)
	}

	response := &pb.RetrieveDataResponse{
		Data: pbData,
	}

	return response, nil
}

// ConvertPBMemoryCellToSchema преобразует экземпляр типа pb.MemoryCell в тип schema.MemoryCell
func ConvertPBMemoryCellToSchema(pbCell *pb.MemoryCell) *schema.MemoryCell {
	schemaCell := &schema.MemoryCell{
		InfoCell: &schema.InfoCell{
			ID:          int64(pbCell.Info.Id),
			DataType:    pbCell.Info.DataType,
			DataSize:    pbCell.Info.DataSize,
			Description: pbCell.Info.Description,
			OwnerID:     int64(pbCell.Info.OwnerId),
		},
		ID:            int64(pbCell.Id),
		InfoID:        int64(pbCell.Info.Id),
		Encrypted:     pbCell.Encrypted,
		KeyValuePairs: pbCell.KeyValuePairs,
		BinaryData:    pbCell.BinaryData,
		FileName:      pbCell.FileName,
	}

	return schemaCell
}

// ConvertPBInfoCellToSchema преобразует экземпляр типа pb.InfoCell в тип schema.InfoCell
func ConvertPBInfoCellToSchema(pbCell *pb.InfoCell) *schema.InfoCell {
	schemaCell := &schema.InfoCell{
		ID:          int64(pbCell.Id),
		DataType:    pbCell.DataType,
		DataSize:    pbCell.DataSize,
		Description: pbCell.Description,
		OwnerID:     int64(pbCell.OwnerId),
	}

	return schemaCell
}

// tokenInterceptor - перехватчик проверяет наличие и валидность токена
func (h *HandlerService) tokenInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Получаем название метода
	methodName := filepath.Base(info.FullMethod)

	// Исключаем методы "Register" и "Authenticate" из проверки токена
	excludedMethods := []string{"Register", "Authenticate", "Authorize"}
	if slices.Contains(excludedMethods, methodName) {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "Failed to get metadata from context")
	}

	// Проверка наличия токена в метаданных
	values := md.Get("token")
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "Token is missing")
	}
	token := values[0]
	if ok, _ := h.gophKeeper.CheckToken(token); !ok {
		return nil, status.Error(codes.Unauthenticated, "Token is invalid")
	}
	return handler(ctx, req)
}

// GetTokenFromContext - получает токен из контекста.
// Возвращает токен первым аргументом, и успех операции вторым.
func GetTokenFromContext(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	tokens := md.Get("token")
	if len(tokens) == 0 {
		return "", false
	}

	return tokens[0], true
}

// ConvertSchemaMemoryCellToPB преобразует экземпляр типа schema.MemoryCell в тип pb.MemoryCell
func ConvertSchemaMemoryCellToPB(schemaCell *schema.MemoryCell) *pb.MemoryCell {
	pbCell := &pb.MemoryCell{
		Id: schemaCell.ID,
		Info: &pb.InfoCell{
			Id:          schemaCell.InfoCell.ID,
			DataType:    schemaCell.InfoCell.DataType,
			DataSize:    schemaCell.InfoCell.DataSize,
			Description: schemaCell.InfoCell.Description,
			OwnerId:     schemaCell.InfoCell.OwnerID,
		},
		Encrypted:     schemaCell.Encrypted,
		KeyValuePairs: schemaCell.KeyValuePairs,
		BinaryData:    schemaCell.BinaryData,
		FileName:      schemaCell.FileName,
	}

	return pbCell
}

// ConvertSchemaInfoCellToPB преобразует экземпляр типа schema.InfoCell в тип pb.InfoCell
func ConvertSchemaInfoCellToPB(schemaCell *schema.InfoCell) *pb.InfoCell {
	pbCell := &pb.InfoCell{
		Id:          schemaCell.ID,
		DataType:    schemaCell.DataType,
		DataSize:    schemaCell.DataSize,
		Description: schemaCell.Description,
		OwnerId:     schemaCell.OwnerID,
	}

	return pbCell
}
