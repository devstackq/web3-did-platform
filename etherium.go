package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"math"
	"math/big"
	"strings"

	"net/http"
)

const trxABI = `[
    {
        "inputs": [
            {
                "internalType": "address payable",
                "name": "_receiver",
                "type": "address"
            }
        ],
        "name": "sendEth",
        "outputs": [],
        "stateMutability": "payable",
        "type": "function"
    },
    {
        "inputs": [],
        "name": "getTxHistory",
        "outputs": [
            {
                "components": [
                    {
                        "internalType": "address",
                        "name": "sender",
                        "type": "address"
                    },
                    {
                        "internalType": "address",
                        "name": "receiver",
                        "type": "address"
                    },
                    {
                        "internalType": "uint256",
                        "name": "amount",
                        "type": "uint256"
                    },
                    {
                        "internalType": "uint256",
                        "name": "timestamp",
                        "type": "uint256"
                    }
                ],
                "internalType": "struct TransactionManager.Transaction[]",
                "name": "",
                "type": "tuple[]"
            }
        ],
        "stateMutability": "view",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "address",
                "name": "_address",
                "type": "address"
            }
        ],
        "name": "getTransactionHistoryByAddress",
        "outputs": [
            {
                "components": [
                    {
                        "internalType": "address",
                        "name": "sender",
                        "type": "address"
                    },
                    {
                        "internalType": "address",
                        "name": "receiver",
                        "type": "address"
                    },
                    {
                        "internalType": "uint256",
                        "name": "amount",
                        "type": "uint256"
                    },
                    {
                        "internalType": "uint256",
                        "name": "timestamp",
                        "type": "uint256"
                    }
                ],
                "internalType": "struct TransactionManager.Transaction[]",
                "name": "",
                "type": "tuple[]"
            }
        ],
        "stateMutability": "view",
        "type": "function"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "internalType": "address",
                "name": "sender",
                "type": "address"
            },
            {
                "indexed": true,
                "internalType": "address",
                "name": "receiver",
                "type": "address"
            },
            {
                "indexed": false,
                "internalType": "uint256",
                "name": "amount",
                "type": "uint256"
            },
            {
                "indexed": false,
                "internalType": "uint256",
                "name": "timestamp",
                "type": "uint256"
            }
        ],
        "name": "TransactionSent",
        "type": "event"
    },
	{
		"inputs": [],
		"name": "getTxHistory",
		"outputs": [
			{
				"components": [
					{
						"internalType": "address",
						"name": "sender",
						"type": "address"
					},
					{
						"internalType": "address",
						"name": "receiver",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "timestamp",
						"type": "uint256"
					}
				],
				"internalType": "struct TransactionManager.Transaction[]",
				"name": "",
				"type": "tuple[]"
			}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`

const contractAddress = "0xe37CC42ea6b89BFCe2E6257FdA9dc04d5FE5960b"

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
		return
	}

	privateKey, err := crypto.HexToECDSA(req.PrivateKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	toAddress := common.HexToAddress(req.RecipientAddress)

	parsedABI, err := abi.JSON(strings.NewReader(trxABI))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amount := big.NewInt(req.Amount)

	if err = e.send(c, parsedABI, fromAddress, toAddress, privateKey, amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error send": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (e *Eth) getTransactionHistory(c *gin.Context) {

	parsedABI, err := abi.JSON(strings.NewReader(trxABI))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encodedData, err := parsedABI.Pack("getTxHistory")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address := common.HexToAddress(contractAddress)

	//read only from smart contract
	result, err := e.cl.CallContract(c, ethereum.CallMsg{
		To:   &address,
		Gas:  0,
		Data: encodedData,
	}, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var transactions []struct {
		Sender    common.Address
		Receiver  common.Address
		Amount    *big.Int
		Timestamp *big.Int
	}

	if err = parsedABI.UnpackIntoInterface(&transactions, "getTxHistory", result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, tx := range transactions {
		fmt.Printf("Sender: %s, Receiver: %s, Amount: %s, Timestamp: %s\n",
			tx.Sender.Hex(), tx.Receiver.Hex(), tx.Amount.String(), tx.Timestamp.String())
	}

	c.Status(http.StatusOK)

}

func (e *Eth) send(ctx context.Context, parsedABI abi.ABI, fromAddress, toAddress common.Address, privateKey *ecdsa.PrivateKey, amount *big.Int) error {

	// Получение nonce
	nonce, err := e.cl.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return err
	}

	// Получение рекомендованной комиссии (gas price)
	gasPrice, err := e.cl.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}

	gasLimit := uint64(300000) // Стандартный лимит для простой транзакции
	// Количество ETH для отправки (в wei)

	data, err := parsedABI.Pack("sendEth", toAddress)
	if err != nil {
		return fmt.Errorf("parsedABI.Pack() failed: %v", err)
	}

	contracAddr := common.HexToAddress(contractAddress)

	// create trx and change state in Blockchain
	tx := types.NewTransaction(nonce, contracAddr, amount, gasLimit, gasPrice, data)

	// Получение ID сети (chain ID)
	chainID, err := e.cl.NetworkID(ctx)
	if err != nil {
		return err
	}
	//chainID := big.NewInt(11155111) // Chain ID для Sepolia
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}

	if err = e.cl.SendTransaction(ctx, signedTx); err != nil {
		return err
	}

	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
	return nil
}
