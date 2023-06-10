package schema

// User представляет структуру данных пользователя
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// InfoCell представляет структуру данных информационной ячейки
type InfoCell struct {
	ID          int64  `json:"id"`
	DataType    string `json:"dataType"`
	DataSize    int32  `json:"dataSize"`
	Description string `json:"description"`
	OwnerID     int64  `json:"ownerId"`
}

// MemoryCell представляет структуру данных ячейки памяти
type MemoryCell struct {
	InfoCell      *InfoCell
	ID            int64             `json:"id"`
	InfoID        int64             `json:"infoId"`
	Encrypted     bool              `json:"encrypted"`
	KeyValuePairs map[string]string `json:"keyValuePairs"`
	BinaryData    []byte            `json:"binaryData"`
	FileName      string            `json:"fileName"`
}
