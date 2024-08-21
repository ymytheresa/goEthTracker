package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/ymytheresa/erc20-token-tracker/contracts"
	"github.com/ymytheresa/erc20-token-tracker/internal/config"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to Ganache
	client, err := ethclient.Dial(cfg.GanacheURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ganache: %v", err)
	}

	// Create a new transactor
	privateKey, err := crypto.HexToECDSA(cfg.DeployerPrivateKey)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337)) // 1337 is the default chain ID for Ganache
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}

	// Deploy the contract
	address, tx, _, err := contracts.DeployMyToken(auth, client, "MyToken", "MTK", big.NewInt(1000000)) // 1,000,000 tokens
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}

	fmt.Printf("Contract deployed to: %s\n", address.Hex())
	fmt.Printf("Transaction hash: %s\n", tx.Hash().Hex())

	// Wait for the transaction to be mined
	_, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatalf("Failed to wait for transaction to be mined: %v", err)
	}

	fmt.Println("Contract deployment confirmed")
}
