package keeper

import (
	"context"
	"errors"
	"fmt"

	"github.com/bubu256/gophkeeper_pet/config"
	"github.com/bubu256/gophkeeper_pet/internal/schema"
	_ "github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Keeper представляет интерфейс для взаимодействия с базой данных.
type Keeper interface {
	GetUserDataInfo(userID int64) ([]*schema.InfoCell, error)
	GetDataByInfoIDs(infoIDs []int64) ([]*schema.MemoryCell, error)
	UpdateMemoryCell(memoryCell schema.MemoryCell) (bool, error)
	UpdateInfoCell(infoCell schema.InfoCell) (bool, error)
	AddData(infoCell schema.InfoCell, memoryCell *schema.MemoryCell) (int64, error)
	DeleteData(infoIDs []int64) (bool, error)
	GetUserByUsername(username string) (*schema.User, error)
	CreateUser(user *schema.User) error
	Ping() error
}

// StoragePG представляет хранилище данных PostgreSQL.
type StoragePG struct {
	db  *pgxpool.Pool
	cfg config.ServerConfig
}

var _ Keeper = &StoragePG{}

// New создает новый экземпляр StoragePG и устанавливает соединение с базой данных.
func New(cfg config.ServerConfig) (*StoragePG, error) {
	dbConfig, err := pgxpool.ParseConfig(cfg.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database DSN: %w", err)
	}

	dbConfig.ConnConfig.Logger = discardLogger{} // Отключение вывода логов pgx

	db, err := pgxpool.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	storage := &StoragePG{
		db:  db,
		cfg: cfg,
	}

	return storage, nil
}

// GetUserByUsername возвращает пользователя по его имени пользователя (username).
func (s *StoragePG) GetUserByUsername(username string) (*schema.User, error) {
	query := `
			SELECT id, username, password_hash
			FROM users
			WHERE username = $1
		`

	row := s.db.QueryRow(context.Background(), query, username)

	user := &schema.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
	)
	if err != nil {
		// if err == pgx.ErrNoRows {
		// 	return user, nil
		// }
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return user, nil
}

// CreateUser создает нового пользователя.
func (s *StoragePG) CreateUser(user *schema.User) error {
	query := `
			INSERT INTO users (username, password_hash)
			VALUES ($1, $2)
			RETURNING id
		`

	err := s.db.QueryRow(
		context.Background(),
		query,
		user.Username,
		user.Password,
	).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}

// GetUserSavedDataInfo возвращает информацию о всех сохраненных данных пользователя по его ID.
func (s *StoragePG) GetUserDataInfo(userID int64) ([]*schema.InfoCell, error) {
	query := `
			SELECT id, data_type, data_size, description, owner_id
			FROM info_cells
			WHERE owner_id = $1
		`

	rows, err := s.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var infoCells []*schema.InfoCell

	for rows.Next() {
		infoCell := &schema.InfoCell{}
		err := rows.Scan(
			&infoCell.ID,
			&infoCell.DataType,
			&infoCell.DataSize,
			&infoCell.Description,
			&infoCell.OwnerID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		infoCells = append(infoCells, infoCell)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return infoCells, nil
}

// GetDataByInfoIDs возвращает данные, соответствующие заданным InfoID.
func (s *StoragePG) GetDataByInfoIDs(infoIDs []int64) ([]*schema.MemoryCell, error) {
	memoryCells := make([]*schema.MemoryCell, 0, len(infoIDs))

	if len(infoIDs) == 0 {
		return memoryCells, nil
	}

	query := `
			SELECT m.id, m.info_id, m.encrypted, m.key_value_pairs, m.binary_data, m.file_name,
				i.data_type, i.data_size, i.description, i.owner_id
			FROM memory_cells m
			INNER JOIN info_cells i ON m.info_id = i.id
			WHERE i.id = ANY($1)
		`

	rows, err := s.db.Query(context.Background(), query, infoIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		memoryCell := &schema.MemoryCell{}
		infoCell := &schema.InfoCell{}

		err := rows.Scan(
			&memoryCell.ID,
			&memoryCell.InfoID,
			&memoryCell.Encrypted,
			&memoryCell.KeyValuePairs,
			&memoryCell.BinaryData,
			&memoryCell.FileName,
			&infoCell.DataType,
			&infoCell.DataSize,
			&infoCell.Description,
			&infoCell.OwnerID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		memoryCell.InfoCell = infoCell
		memoryCells = append(memoryCells, memoryCell)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return memoryCells, nil
}

// UpdateMemoryCell обновляет данные ячейки памяти на основе InfoID.
func (s *StoragePG) UpdateMemoryCell(memoryCell schema.MemoryCell) (bool, error) {
	query := `
			UPDATE memory_cells
			SET encrypted = $1, key_value_pairs = $2, binary_data = $3, file_name = $4
			WHERE info_id = $5
		`

	result, err := s.db.Exec(
		context.Background(),
		query,
		memoryCell.Encrypted,
		memoryCell.KeyValuePairs,
		memoryCell.BinaryData,
		memoryCell.FileName,
		memoryCell.InfoID,
	)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return false, errors.New("no rows affected")
	}

	return true, nil
}

// UpdateInfoCell обновляет данные информационной ячейки на основе ID.
func (s *StoragePG) UpdateInfoCell(infoCell schema.InfoCell) (bool, error) {
	query := `
			UPDATE info_cells
			SET data_type = $1, data_size = $2, description = $3
			WHERE id = $4
		`

	result, err := s.db.Exec(
		context.Background(),
		query,
		infoCell.DataType,
		infoCell.DataSize,
		infoCell.Description,
		infoCell.ID,
	)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return false, errors.New("no rows affected")
	}

	return true, nil
}

// AddData добавляет новые данные в базу данных.
func (s *StoragePG) AddData(infoCell schema.InfoCell, memoryCell *schema.MemoryCell) (int64, error) {
	insertQuery := `
		WITH inserted_info AS (
			INSERT INTO info_cells (data_type, data_size, description, owner_id)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		)
		INSERT INTO memory_cells (info_id, encrypted, key_value_pairs, binary_data, file_name)
		SELECT id, $5, $6, $7, $8
		FROM inserted_info
		RETURNING info_id
	`
	// log.Printf("infoCell.OwnerID %v", infoCell.OwnerID)
	var infoID int64
	err := s.db.QueryRow(
		context.Background(),
		insertQuery,
		infoCell.DataType,
		infoCell.DataSize,
		infoCell.Description,
		infoCell.OwnerID,
		memoryCell.Encrypted,
		memoryCell.KeyValuePairs,
		memoryCell.BinaryData,
		memoryCell.FileName,
	).Scan(&infoID)
	if err != nil {
		return 0, fmt.Errorf("failed to execute insert query: %w", err)
	}

	return infoID, nil
}

// DeleteData удаляет данные из базы данных на основе заданных InfoID.
func (s *StoragePG) DeleteData(infoIDs []int64) (bool, error) {
	query := `
			DELETE FROM info_cells
			WHERE id = ANY($1)
		`

	result, err := s.db.Exec(context.Background(), query, infoIDs)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return false, errors.New("no rows affected")
	}

	return true, nil
}

// Ping проверяет доступность соединения с базой данных.
func (s *StoragePG) Ping() error {
	conn, err := s.db.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	// var pingResult pgconn.PingResult
	err = conn.Conn().Ping(context.Background())
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	return nil
}

// discardLogger является заглушкой для логгера pgx.
type discardLogger struct{}

// Log реализует метод Log интерфейса pgx.Logger.
func (d discardLogger) Log(_ context.Context, _ pgx.LogLevel, _ string, _ map[string]interface{}) {}
