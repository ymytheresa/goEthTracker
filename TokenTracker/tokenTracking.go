package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// TransferEvent represents a simplified transfer event
type TransferEvent struct {
	From  string
	To    string
	Value string
}

// TokenTracker handles the tracking of ERC20 token events
type TokenTracker struct {
	client          *ethclient.Client
	contractAddress common.Address
	transfers       []TransferEvent
}

// NewTokenTracker creates a new TokenTracker instance
func NewTokenTracker() (*TokenTracker, error) {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Get the Ganache URL and contract address from environment variables
	ganacheURL := os.Getenv("GANACHE_URL")
	contractAddressStr := os.Getenv("CONTRACT_ADDRESS")

	if ganacheURL == "" || contractAddressStr == "" {
		return nil, fmt.Errorf("GANACHE_URL or CONTRACT_ADDRESS not set in .env file")
	}

	// Connect to Ganache
	client, err := ethclient.Dial(ganacheURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ganache: %w", err)
	}

	contractAddress := common.HexToAddress(contractAddressStr)

	return &TokenTracker{
		client:          client,
		contractAddress: contractAddress,
		transfers:       []TransferEvent{},
	}, nil
}

// StartTracking begins listening for Transfer events
func (t *TokenTracker) StartTracking() error {
	// TODO: Implement this function
	// 1. Set up event listening
	// 2. Process incoming events
	// 3. Store events in the transfers slice
	// 4. Print event details to console
	return nil
}

// GetTransfers returns all tracked transfer events
func (t *TokenTracker) GetTransfers() []TransferEvent {
	return t.transfers
}
