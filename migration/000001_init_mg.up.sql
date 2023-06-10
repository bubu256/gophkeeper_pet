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
