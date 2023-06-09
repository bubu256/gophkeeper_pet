// Package config - содержит описание структур конфигурации сервера и клиента и функции для создания конфигурации
package config

// ServerConfig - конфигурация для сервера
type ServerConfig struct {
	Port    string
	Address string
}

// ClientConfig - конфигурация для клиента
type ClientConfig struct {
	ServerAddress string
}
