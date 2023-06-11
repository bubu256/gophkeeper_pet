package goph

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/bubu256/gophkeeper_pet/config"
	"github.com/bubu256/gophkeeper_pet/internal/schema"
	"github.com/bubu256/gophkeeper_pet/pkg/keeper"
)

// Goph представляет интерфейс для работы с бизнес-логикой приложения.
type Goph interface {
	CheckToken(token string) (bool, error)
	CreateUser(username, password string) error
	Authenticate(username, password string) (string, error)
	SaveData(userID int64, memoryCell *schema.MemoryCell) (int64, error)
	GetUserDataInfo(userID int64) ([]*schema.InfoCell, error)
	GetUserMemoryData(userID int64, infoIDs []int64) ([]*schema.MemoryCell, error)
	UserExists(username string) (bool, error)
	GetUserIDFromToken(token string) (int64, error)
}

// GophLogic представляет реализацию интерфейса Goph.
type GophLogic struct {
	secretKey []byte
	keeper    keeper.Keeper
}

var _ Goph = &GophLogic{}

// New создает новый экземпляр GophLogic с заданным Keeper и ServerConfig.
// Генерирует случайный секретный ключ.
func New(keeper keeper.Keeper, config config.ServerConfig) *GophLogic {
	// Генерация случайного секретного ключа
	secretKey := make([]byte, 32)
	_, err := rand.Read(secretKey)
	if err != nil {
		panic(fmt.Errorf("failed to generate secret key: %w", err))
	}

	return &GophLogic{
		secretKey: secretKey,
		keeper:    keeper,
	}
}

// GenerateToken - генерирует токен из ID пользователя и секретного ключа
func (g *GophLogic) GenerateToken(userID int64) (string, error) {
	idBytes := make([]byte, 4)
	if userID < 0 || userID > (1<<31-1) {
		return "", fmt.Errorf("userID is out of range")
	}
	binary.BigEndian.PutUint32(idBytes, uint32(userID))

	h := hmac.New(sha256.New, g.secretKey)
	h.Write(idBytes)
	signature := h.Sum(nil)

	token := append(idBytes, signature...)
	return hex.EncodeToString(token), nil
}

// CheckToken - проверяет токен на подлинность
func (g *GophLogic) CheckToken(token string) (bool, error) {
	decodeToken, err := hex.DecodeString(token)
	if err != nil {
		return false, err
	}
	idUser := decodeToken[:4]
	sing := decodeToken[4:]
	h := hmac.New(sha256.New, g.secretKey)
	h.Write(idUser)
	dst := h.Sum(nil)
	return hmac.Equal(sing, dst), nil
}

// GetUserIDFromToken - получает ID пользователя из токена.
func (g *GophLogic) GetUserIDFromToken(token string) (int64, error) {
	decodedToken, err := hex.DecodeString(token)
	if err != nil {
		return 0, err
	}

	if len(decodedToken) < 4 {
		return 0, fmt.Errorf("invalid token length")
	}

	idBytes := decodedToken[:4]
	userID := binary.BigEndian.Uint32(idBytes)

	return int64(userID), nil
}

// CreateUser создает нового пользователя.
// Хеширует логин и вызывает метод Keeper для сохранения пользователя.
func (g *GophLogic) CreateUser(username, password string) error {
	hashedPassword := hashPassword(password)

	// hashedUsername := hashUsername(username)

	user := &schema.User{
		Username: username,
		Password: hashedPassword,
	}

	err := g.keeper.CreateUser(user)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// ExistUser - проверяет существует ли пользователь
func (g *GophLogic) UserExists(username string) (bool, error) {
	user, err := g.keeper.GetUserByUsername(username)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

// GetUserID - возвращает ID пользователя
func (g *GophLogic) GetUserID(username string) (int64, error) {
	user, err := g.keeper.GetUserByUsername(username)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

// Authenticate выполняет аутентификацию пользователя.
// Хеширует пароль и проверяет его с хранимым хешем в Keeper.
// Возвращает токен или ошибку, если аутентификация не удалась.
func (g *GophLogic) Authenticate(username, password string) (string, error) {
	user, err := g.keeper.GetUserByUsername(username)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user: %w", err)
	}

	if hashPassword(password) != user.Password {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	return g.GenerateToken(user.ID)
}

// SaveData сохраняет новые данные для пользователя.
// Вызывает соответствующий метод Keeper для записи данных в базу данных.
func (g *GophLogic) SaveData(userID int64, memoryCell *schema.MemoryCell) (int64, error) {
	memoryCell.InfoCell.OwnerID = userID
	// log.Printf("infoCell: %+v", memoryCell.InfoCell)
	infoID, err := g.keeper.AddData(*memoryCell.InfoCell, memoryCell)
	if err != nil {
		return infoID, fmt.Errorf("failed to save memory cell: %w", err)
	}

	return infoID, nil
}

// GetUserData возвращает информацию о данных пользователя.
// Вызывает метод Keeper для получения данных из базы данных.
func (g *GophLogic) GetUserDataInfo(userID int64) ([]*schema.InfoCell, error) {
	infoCells, err := g.keeper.GetUserDataInfo(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve info cells: %w", err)
	}

	return infoCells, nil
}

// GetUserMemoryData возвращает данные пользователя для указанных идентификаторов InfoID.
// Проверяет, принадлежат ли идентификаторы пользователю, и вызывает соответствующий метод Keeper для получения данных из базы данных.
func (g *GophLogic) GetUserMemoryData(userID int64, infoIDs []int64) ([]*schema.MemoryCell, error) {
	infoCells, err := g.keeper.GetUserDataInfo(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user data: %w", err)
	}

	// Создаем слайс из InfoCell.ID
	infoCellIDs := make([]int64, 0, len(infoCells))
	for _, infoCell := range infoCells {
		infoCellIDs = append(infoCellIDs, infoCell.ID)
	}

	// Фильтруем infoIDs, оставляем только те ID, которые есть в infoCellIDs
	filteredInfoIDs := make([]int64, 0, len(infoIDs))
	for _, id := range infoIDs {
		for _, infoCellID := range infoCellIDs {
			if id == infoCellID {
				filteredInfoIDs = append(filteredInfoIDs, id)
				break
			}
		}
	}

	// Получаем данные MemoryCell по отфильтрованным ID
	memoryCells, err := g.keeper.GetDataByInfoIDs(filteredInfoIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user memory cells: %w", err)
	}

	return memoryCells, nil
}

// hashUsername хеширует логин пользователя.
func hashPassword(username string) string {
	hashedUsername := sha256.Sum256([]byte(username))
	return base64.URLEncoding.EncodeToString(hashedUsername[:])
}
