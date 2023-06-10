-- Файл миграции для создания таблиц данных

CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS info_cells (
  id SERIAL PRIMARY KEY,
  data_type VARCHAR(255) NOT NULL,
  data_size INT NOT NULL,
  description TEXT,
  owner_id INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (owner_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS memory_cells (
  id SERIAL PRIMARY KEY,
  info_id INT NOT NULL,
  encrypted BOOLEAN NOT NULL,
  key_value_pairs JSONB,
  binary_data BYTEA,
  file_name VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (info_id) REFERENCES info_cells (id)
);

-- Вставка тестовых данных в таблицу пользователей
INSERT INTO users (username, password_hash)
VALUES
  ('john.doe', 'password_hash_1'),
  ('jane.smith', 'password_hash_2');

-- Вставка тестовых данных в таблицу информационных ячеек
INSERT INTO info_cells (data_type, data_size, description, owner_id)
VALUES
  ('data_type_1', 1024, 'Description 1', 1),
  ('data_type_2', 2048, 'Description 2', 2);

-- Вставка тестовых данных в таблицу ячеек памяти
INSERT INTO memory_cells (info_id, encrypted, key_value_pairs, binary_data, file_name)
VALUES
  (1, FALSE, '{"key1": "value1", "key2": "value2"}', E'\\x0123456789ABCDEF', 'file1.txt'),
  (2, TRUE, '{"key3": "value3", "key4": "value4"}', E'\\xAABBCCDDEEFF', 'file2.txt');
