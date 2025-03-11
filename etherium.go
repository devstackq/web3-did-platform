package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"math"
	"math/big"

	"net/http"
)

type Eth struct {
	cl ethclient.Client
}

type txRequest struct {
	PrivateKey       string `json:"private_key"`
	RecipientAddress string `json:"recipient_address"`
	Amount           int64  `json:"amount"`
}

func New(projectID string) (*Eth, error) {
	// Подключение к Ethereum через Infura
	//url := fmt.Sprint("https://mainnet.infura.io/v3/", projectID)//eth prod env
	url := fmt.Sprint("https://sepolia.infura.io/v3/", projectID)

	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Eth{cl: *client}, nil
}

func (e *Eth) getBalance(c *gin.Context) {
	address := c.Param("address")

	if ok := common.IsHexAddress(address); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "is not HEX address"})
		return
	}

	balance, err := e.cl.BalanceAt(c, common.HexToAddress(address), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	balanceETH := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(math.Pow10(18)))
	fmt.Println(balanceETH.String(), "balanceETH")

	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"balance": balanceETH.String(),
	})
}

func (e *Eth) sendTransaction(c *gin.Context) {

	var req txRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	privateKey, err := crypto.HexToECDSA(req.PrivateKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid private key"})
		return
	}

	// Адрес отправителя
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cast public key to ECDSA"})
		return
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Адрес получателя
	toAddress := common.HexToAddress(req.RecipientAddress)

	// Получение nonce
	nonce, err := e.cl.PendingNonceAt(c, fromAddress)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get nonce"})
		return
	}

	// Получение рекомендованной комиссии (gas price)
	gasPrice, err := e.cl.SuggestGasPrice(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get gas price"})
		return
	}

	gasLimit := uint64(21000) // Стандартный лимит для простой транзакции
	//gasPrice := big.NewInt(20000000000) // 20 Gwei
	// Количество ETH для отправки (в wei)
	amount := big.NewInt(req.Amount)

	fmt.Println(gasPrice, "gasPrice")

	// Создание транзакции
	tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)

	// Получение ID сети (chain ID)
	chainID, err := e.cl.NetworkID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chain ID"})
		return
	}
	//chainID := big.NewInt(11155111) // Chain ID для Sepolia
	// Подписание транзакции
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to sign transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	// Отправка транзакции
	if err = e.cl.SendTransaction(c, signedTx); err != nil {
		errMsg := fmt.Sprintf("Failed to send transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	// Возвращаем хэш транзакции
	c.JSON(http.StatusOK, gin.H{
		"transaction_hash": signedTx.Hash().Hex(),
	})
}
