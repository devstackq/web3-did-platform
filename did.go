package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func createNewDID(c *gin.Context) {
	// Генерация нового приватного ключа
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate private key"})
		return
	}

	// Получение публичного ключа из приватного
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cast public key to ECDSA"})
		return
	}

	// Генерация Ethereum-адреса из публичного ключа
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	// Формирование DID
	did := fmt.Sprintf("did:ethr:%s", address)

	// Возвращаем DID и приватный ключ (для демонстрации, в реальном приложении приватный ключ не должен передаваться клиенту)
	c.JSON(http.StatusOK, gin.H{
		"did":         did,
		"address":     address,
		"private_key": fmt.Sprintf("%x", crypto.FromECDSA(privateKey)), // Приватный ключ в hex-формате
	})
}

// createNewDID генерирует DID на основе Ethereum-адреса
func createDID(c *gin.Context) {

	address := c.Query("address")

	// Формирование DID
	did := fmt.Sprintf("did:ethr:%s", address)

	// Возвращаем DID и приватный ключ (для демонстрации, в реальном приложении приватный ключ не должен передаваться клиенту)
	c.JSON(http.StatusOK, gin.H{
		"did":     did,
		"address": address,
	})
}
