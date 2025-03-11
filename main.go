package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Web3 Platform",
		})
	})

	ethCl, err := New("5d4822d03f3940748e09a54592655fd5")
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/did/create", createDID)

	r.GET("/eth/balance/:address", ethCl.getBalance)
	r.POST("/eth/send", ethCl.sendTransaction)
	r.POST("/eth/history", ethCl.getTransactionHistory)

	if err = r.Run(); err != nil {
		log.Fatal(err)
	} // 8080
}

// TODO solidity add
