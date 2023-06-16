// Package storage - позволяет локально хранить информацию о данных пользователя.
package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/bubu256/gophkeeper_pet/internal/proto/pb"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

// Storage представляет структуру для хранения в памяти набора MemoryCell.
type Storage struct {
	data         []*pb.MemoryCell
	dumpFilePath string
}

// NewStorage - возвращает экземпляр хранилища
func NewStorage() *Storage {
	dirPath, err := CreateAppDirectory("gophkeeper")
	if err != nil {
		log.Printf("failed application directory creation %v", err)
	}
	dumpFilePath := filepath.Join(dirPath, ".dump")
	return &Storage{
		dumpFilePath: dumpFilePath,
	}
}

// Sync - сохраняет данные в памяти.
func (s *Storage) Sync(cells []*pb.MemoryCell) {
	// s.data = append(s.data, cells...)
	s.data = cells
}

// Получить информацию: возвращает слайс pb.InfoCell, сформированный из набора MemoryCell.
func (s *Storage) GetInformation() []*pb.InfoCell {
	var infoCells []*pb.InfoCell
	for _, cell := range s.data {
		infoCells = append(infoCells, cell.Info)
	}
	return infoCells
}

// Получить данные по ID: возвращает pb.MemoryCell по указанному InfoID.
func (s *Storage) GetDataByID(infoID int64) (*pb.MemoryCell, error) {
	for _, cell := range s.data {
		if cell.Info.Id == infoID {
			return cell, nil
		}
	}
	return nil, errors.New("data not found")
}

// Получить данные по ID: возвращает pb.MemoryCell по указанному InfoID.
func (s *Storage) GetData() []*pb.MemoryCell {
	return s.data
}

// Сделать дамп на диск: шифрует данные и сохраняет их в файл на диске.
func (s *Storage) DumpToFile(password string) error {
	// Преобразование данных в бинарный формат
	data, err := proto.Marshal(&pb.RetrieveDataResponse{Data: s.data})
	if err != nil {
		return errors.Wrap(err, "failed to marshal data to binary")
	}

	// Генерация ключа из пароля
	key := GenerateKeyFromPassword(password)

	// Шифрование данных
	encryptedData, err := Encrypt(data, key)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt data")
	}

	// создание файла
	file, err := os.Create(s.dumpFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to create dump file")
	}
	defer file.Close()

	// Сохранение данных в файл
	err = ioutil.WriteFile(s.dumpFilePath, encryptedData, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write data to file")
	}

	return nil
}

// Загрузить из дампа на диске: считывает и дешифрует данные из файла на диске и сохраняет их в структуре.
func (s *Storage) LoadFromDump(password string) error {
	// Чтение данных из файла
	encryptedData, err := ioutil.ReadFile(s.dumpFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to read data from file")
	}

	// Генерация ключа из пароля
	key := GenerateKeyFromPassword(password)

	// Расшифровка данных
	decryptedData, err := Decrypt(encryptedData, key)
	if err != nil {
		return errors.Wrap(err, "failed to decrypt data")
	}

	// Преобразование данных из бинарного формата
	var response pb.RetrieveDataResponse
	err = proto.Unmarshal(decryptedData, &response)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal data from binary")
	}

	// Обновление данных в структуре
	s.data = response.Data

	return nil
}

// Функция для генерации ключа из пароля
func GenerateKeyFromPassword(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:32]
}

// Encrypt шифрует данные с использованием AES и указанного ключа.
func Encrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	ciphertext = append(nonce, ciphertext...)

	return ciphertext, nil
}

// Decrypt дешифрует данные с использованием AES и указанного ключа.
func Decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("invalid ciphertext")
	}

	nonce := ciphertext[:12]
	encryptedData := ciphertext[12:]

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// CreateAppDirectory - создает директорию приложения для хранения дампа. Возвращает путь к будущему дампу.
func CreateAppDirectory(dirName string) (string, error) {
	// Получаем путь к директории пользователя
	userDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Создаем папку в директории пользователя
	dirPath := filepath.Join(userDir, dirName)
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	log.Println(dirPath)

	// Возвращаем путь к созданной папке
	return dirPath, nil
}
