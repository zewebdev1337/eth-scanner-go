package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	apiKey := os.Getenv("INFURA_PROJECT_ID")
	if apiKey == "" {
		log.Fatalf("The 'INFURA_PROJECT_ID' environment variable is not set.")
	}

	infuraURL := fmt.Sprintf("wss://mainnet.infura.io/ws/v3/%s", apiKey)

	// Connect to Ethereum network
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum network: %v", err)
	}
	defer client.Close()

	// Check connection
	_, err = client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum network: %v", err)
	}
	fmt.Println("Connected to Ethereum network!")

	// Subscribe to new blocks
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatalf("Failed to subscribe to new headers: %v", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			handleBlock(client, header.Number)
		}
	}
}

func handleBlock(client *ethclient.Client, blockNumber *big.Int) {
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Printf("Failed to get block %v: %v", blockNumber, err)
		return
	}

	for _, tx := range block.Transactions() {
		if tx.To() == nil {
			txHash := tx.Hash().Hex()
			etherscanLink := fmt.Sprintf("https://etherscan.io/tx/%s", txHash)
			fmt.Printf("New contract created with hash: [%s](%s)\n", txHash, etherscanLink)
		}
	}
}
