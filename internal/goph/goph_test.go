package goph_test

import (
	"testing"

	"github.com/bubu256/gophkeeper_pet/config"
	"github.com/bubu256/gophkeeper_pet/internal/goph"
	"github.com/stretchr/testify/assert"
)

func TestCheckToken_ValidToken(t *testing.T) {
	secretKey := []byte("secret_key")
	gophLogic := goph.New(nil, config.ServerConfig{})
	gophLogic.SetSecretKey(secretKey)

	userID := int64(123)
	token, _ := gophLogic.GenerateToken(userID)

	valid, err := gophLogic.CheckToken(token)
	assert.True(t, valid)
	assert.NoError(t, err)
}

func TestCheckToken_InvalidToken(t *testing.T) {
	secretKey := []byte("secret_key")
	gophLogic := goph.New(nil, config.ServerConfig{})
	gophLogic.SetSecretKey(secretKey)

	// Create an invalid token by modifying a valid token
	userID := int64(123)
	validToken, _ := gophLogic.GenerateToken(userID)
	invalidToken := validToken[1:] + "f"

	valid, err := gophLogic.CheckToken(invalidToken)
	assert.False(t, valid)
	assert.NoError(t, err)
}

func TestGenerateToken(t *testing.T) {
	secretKey := []byte("secret_key")
	gophLogic := goph.New(nil, config.ServerConfig{})
	gophLogic.SetSecretKey(secretKey)

	userID := int64(123)
	token, err := gophLogic.GenerateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestHashPassword(t *testing.T) {
	username := "test_user"
	expectedHash := "EWATCHX9oIEsmcXj8aA1FkcaY3DE-XEpsiGTjrR2PmM="

	hash := goph.HashPassword(username)
	assert.Equal(t, expectedHash, hash)
}
