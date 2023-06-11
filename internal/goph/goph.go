package goph

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/bubu256/gophkeeper_pet/config"
	"github.com/bubu256/gophkeeper_pet/internal/schema"
	"github.com/bubu256/gophkeeper_pet/pkg/keeper"
)

// Goph представляет интерфейс для работы с бизнес-логикой приложения.
type Goph interface {
	GenerateToken(userID int64) ([]byte, error)
	CheckToken(token []byte) (bool, error)
	CreateUser(username, password string) error
	Authenticate(username, password string) ([]byte, error)
	SaveData(userID int64, memoryCell schema.MemoryCell) error
	GetUserDataInfo(userID int64) ([]*schema.InfoCell, error)
	GetUserMemoryData(userID int64, infoIDs []int64) ([]*schema.MemoryCell, error)
}

// GophLogic представляет реализацию интерфейса Goph.
type GophLogic struct {
	secretKey []byte
	keeper    keeper.Keeper
}

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

// GenerateToken генерирует новый токен на основе ID пользователя.
// В первых двух байтах токена зашифрован ID пользователя, а остальная часть получена шифрованием ID с использованием секретного ключа.
func (g *GophLogic) GenerateToken(userID int64) ([]byte, error) {
	// Шифрование ID пользователя
	idBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(idBytes, uint64(userID))

	block, err := aes.NewCipher(g.secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Зашифрование ID
	ciphertext := make([]byte, aes.BlockSize+len(idBytes))
	iv := ciphertext[:aes.BlockSize]
	_, err = rand.Read(iv)
	if err != nil {
		return nil, fmt.Errorf("failed to generate IV: %w", err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], idBytes)

	return ciphertext, nil
}

// CheckToken проверяет корректность токена.
// Извлекает ID пользователя из токена и проверяет, что остаток токена соответствует шифрованному ID.
func (g *GophLogic) CheckToken(token []byte) (bool, error) {
	if len(token) < 2+aes.BlockSize {
		return false, errors.New("invalid token length")
	}

	block, err := aes.NewCipher(g.secretKey)
	if err != nil {
		return false, fmt.Errorf("failed to create cipher block: %w", err)
	}

	iv := token[2 : 2+aes.BlockSize]
	encryptedID := token[2+aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	decryptedID := make([]byte, len(encryptedID))
	stream.XORKeyStream(decryptedID, encryptedID)

	id := int64(binary.LittleEndian.Uint64(decryptedID))
	rest := token[2+aes.BlockSize:]
	encryptedRest := make([]byte, len(rest))
	stream.XORKeyStream(encryptedRest, rest)

	return id == int64(binary.LittleEndian.Uint64(encryptedRest)), nil
}

// CreateUser создает нового пользователя.
// Хеширует логин и вызывает метод Keeper для сохранения пользователя.
func (g *GophLogic) CreateUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	hashedUsername := hashUsername(username)

	user := &schema.User{
		Username: hashedUsername,
		Password: string(hashedPassword),
	}

	err = g.keeper.CreateUser(user)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// Authenticate выполняет аутентификацию пользователя.
// Хеширует пароль и проверяет его с хранимым хешем в Keeper.
// Возвращает токен или ошибку, если аутентификация не удалась.
func (g *GophLogic) Authenticate(username, password string) ([]byte, error) {
	hashedUsername := hashUsername(username)

	user, err := g.keeper.GetUserByUsername(hashedUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return g.GenerateToken(user.ID)
}

// SaveData сохраняет новые данные для пользователя.
// Вызывает соответствующий метод Keeper для записи данных в базу данных.
func (g *GophLogic) SaveData(userID int64, memoryCell schema.MemoryCell) error {
	memoryCell.InfoCell.OwnerID = userID
	err := g.keeper.AddData(*memoryCell.InfoCell, memoryCell)
	if err != nil {
		return fmt.Errorf("failed to save memory cell: %w", err)
	}

	return nil
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
func hashUsername(username string) string {
	hashedUsername := sha256.Sum256([]byte(username))
	return string(hashedUsername[:])
}
