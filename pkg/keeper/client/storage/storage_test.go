package storage_test

import (
	"testing"

	"github.com/bubu256/gophkeeper_pet/internal/proto/pb"
	"github.com/bubu256/gophkeeper_pet/pkg/keeper/client/storage"
	"github.com/stretchr/testify/assert"
)

func TestSync(t *testing.T) {
	s := storage.NewStorage()

	// Создаем набор MemoryCell
	cells := []*pb.MemoryCell{
		{
			Info: &pb.InfoCell{
				Id:          1,
				DataType:    "string",
				DataSize:    10,
				Description: "Cell 1",
				OwnerId:     1001,
			},
		},
		{
			Info: &pb.InfoCell{
				Id:          2,
				DataType:    "int",
				DataSize:    4,
				Description: "Cell 2",
				OwnerId:     1002,
			},
		},
	}

	infoexpected := []*pb.InfoCell{cells[0].Info, cells[1].Info}
	// Вызываем метод Sync
	s.Sync(cells)

	// Проверяем, что данные были сохранены в хранилище
	assert.Equal(t, infoexpected, s.GetInformation())
}

func TestGetDataByID(t *testing.T) {
	s := storage.NewStorage()

	// Создаем набор MemoryCell
	cells := []*pb.MemoryCell{
		{
			Info: &pb.InfoCell{
				Id:          1,
				DataType:    "string",
				DataSize:    10,
				Description: "Cell 1",
				OwnerId:     1001,
			},
		},
		{
			Info: &pb.InfoCell{
				Id:          2,
				DataType:    "int",
				DataSize:    4,
				Description: "Cell 2",
				OwnerId:     1002,
			},
		},
	}

	// Вызываем метод Sync
	s.Sync(cells)

	// Проверяем получение данных по ID
	cell, err := s.GetDataByID(2)
	assert.NoError(t, err)
	assert.Equal(t, cells[1], cell)

	// Проверяем получение данных по несуществующему ID
	_, err = s.GetDataByID(3)
	assert.Error(t, err)
}

func TestDumpAndLoadFromDump(t *testing.T) {
	s := storage.NewStorage()

	// Создаем набор MemoryCell
	cells := []*pb.MemoryCell{
		{
			Info: &pb.InfoCell{
				Id:          1,
				DataType:    "string",
				DataSize:    10,
				Description: "Cell 1",
				OwnerId:     1001,
			},
		},
		{
			Info: &pb.InfoCell{
				Id:          2,
				DataType:    "int",
				DataSize:    4,
				Description: "Cell 2",
				OwnerId:     1002,
			},
		},
	}

	// Вызываем метод Sync
	s.Sync(cells)

	// Выполняем дамп данных в файл
	password := "mypassword"
	err := s.DumpToFile(password)
	assert.NoError(t, err)

	// Создаем новый экземпляр хранилища
	newStorage := storage.NewStorage()

	// Загружаем данные из дампа
	err = newStorage.LoadFromDump(password)
	assert.NoError(t, err)

	// Проверяем, что загруженные данные совпадают с исходными данными
	loadCells := newStorage.GetData()
	assert.Len(t, loadCells, len(cells)) // Проверяем, что количество элементов совпадает

	for i, expectedCell := range cells {
		actualCell := loadCells[i]

		assert.Equal(t, expectedCell.Info.Id, actualCell.Info.Id)
		assert.Equal(t, expectedCell.Info.DataType, actualCell.Info.DataType)
		assert.Equal(t, expectedCell.Info.DataSize, actualCell.Info.DataSize)
		assert.Equal(t, expectedCell.Info.Description, actualCell.Info.Description)
		assert.Equal(t, expectedCell.Info.OwnerId, actualCell.Info.OwnerId)
		assert.Equal(t, expectedCell.Encrypted, actualCell.Encrypted)
		assert.Equal(t, expectedCell.KeyValuePairs, actualCell.KeyValuePairs)
		assert.Equal(t, expectedCell.BinaryData, actualCell.BinaryData)
		assert.Equal(t, expectedCell.FileName, actualCell.FileName)
	}

}

func TestGenerateKeyFromPassword(t *testing.T) {
	password := "mypassword"

	// Генерируем ключ из пароля
	key := storage.GenerateKeyFromPassword(password)

	// Проверяем, что ключ сгенерирован
	assert.NotNil(t, key)
	assert.Equal(t, 32, len(key))
}

func TestEncryptAndDecrypt(t *testing.T) {
	data := []byte("hello world")
	password := "mypassword"

	// Генерируем ключ из пароля
	key := storage.GenerateKeyFromPassword(password)

	// Шифруем данные
	encryptedData, err := storage.Encrypt(data, key)
	assert.NoError(t, err)
	assert.NotNil(t, encryptedData)

	// Дешифруем данные
	decryptedData, err := storage.Decrypt(encryptedData, key)
	assert.NoError(t, err)
	assert.NotNil(t, decryptedData)

	// Проверяем, что исходные данные совпадают с расшифрованными данными
	assert.Equal(t, data, decryptedData)
}
