// Package config - содержит описание структур конфигурации сервера и клиента и функции для создания конфигурации
package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

// ServerConfig - конфигурация для сервера
type ServerConfig struct {
	Port        string `env:"SERVER_PORT"`
	Address     string `env:"SERVER_ADDRESS"`
	DatabaseDSN string `env:"DATABASE_DSN"`
}

// ClientConfig - конфигурация для клиента
type ClientConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
}

// LoadFromEnv заполняет конфигурацию сервера из переменных окружения.
func (s *ServerConfig) LoadFromEnv() error {
	err := env.Parse(s)
	if err != nil {
		log.Fatalf("failed to parse server config from environment: %v", err)
		return err
	}
	return nil
}

// LoadFromEnv заполняет конфигурацию клиента из переменных окружения.
func (c *ClientConfig) LoadFromEnv() error {
	err := env.Parse(c)
	if err != nil {
		log.Fatalf("failed to parse client config from environment: %v", err)
		return err
	}
	return nil
}

// GetServerConfig возвращает экземпляр конфигурации сервера.
func GetServerConfig() (ServerConfig, error) {
	cfg := ServerConfig{}
	err := cfg.LoadFromEnv()
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

// GetClientConfig возвращает экземпляр конфигурации клиента.
func GetClientConfig() (ClientConfig, error) {
	cfg := ClientConfig{}
	err := cfg.LoadFromEnv()
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
