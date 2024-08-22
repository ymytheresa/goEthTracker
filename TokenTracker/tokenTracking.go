package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ymytheresa/erc20-token-tracker/ERC20Token/contractsgo"
	"github.com/ymytheresa/erc20-token-tracker/ERC20Token/interact"
)

// TransferEvent represents a simplified transfer event
type TransferEvent struct {
	From  string
	To    string
	Value string
}

// TokenTracker handles the tracking of ERC20 token events
type TokenTracker struct {
	contract            *contractsgo.TestERC20
	client              *ethclient.Client
	contractAddress     common.Address
	transfers           []TransferEvent
	totalTransferredOut *big.Int
}

func setTokenTracker() TokenTracker {
	contractAddress, err := GetContractAddress()
	if err != nil {
		log.Fatal(err)
	}

	return TokenTracker{
		client:              interact.GetClient(),
		contractAddress:     contractAddress,
		transfers:           []TransferEvent{},
		totalTransferredOut: big.NewInt(0),
	}
}

// StartTracking begins listening for Transfer events
func (t *TokenTracker) StartTracking() error {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{t.contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := t.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %v", err)
	}

	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Fatal(err)
			case vLog := <-logs:
				event, err := t.parseTransferEvent(vLog)
				if err != nil {
					log.Printf("Failed to parse transfer event: %v", err)
					continue
				}
				t.transfers = append(t.transfers, event)
				fmt.Printf("Transfer: %s tokens from %s to %s\n", event.Value, event.From, event.To)
			}
		}
	}()

	return nil
}

// parseTransferEvent parses a log entry into a TransferEvent
func (t *TokenTracker) parseTransferEvent(vLog types.Log) (TransferEvent, error) {
	event := TransferEvent{}
	contract, err := contractsgo.NewTestERC20(t.contractAddress, t.client)
	if err != nil {
		return event, fmt.Errorf("failed to create contract instance: %v", err)
	}

	transferEvent, err := contract.ParseTransfer(vLog)
	if err != nil {
		return event, fmt.Errorf("failed to parse transfer event: %v", err)
	}

	event.From = transferEvent.From.Hex()
	event.To = transferEvent.To.Hex()
	event.Value = transferEvent.Value.String()

	return event, nil
}

// GetTransfers returns all tracked transfer events
func (t *TokenTracker) GetTransfers() []TransferEvent {
	return t.transfers
}

func GetContractAddress() (common.Address, error) {
	filePath := "./contract_address.txt"
	addressBytes, err := os.ReadFile(filePath)
	if err != nil {
		return common.Address{}, err
	}
	addressString := string(addressBytes)
	return common.HexToAddress(addressString), nil
}

func (t *TokenTracker) MonitorContractTransfers(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		t.updateTotalTransferredOut()
		fmt.Printf("Total amount transferred from contract: %s\n", t.totalTransferredOut.String())
	}
}

func (t *TokenTracker) updateTotalTransferredOut() {
	if t.contract == nil {
		var err error
		t.contract, err = contractsgo.NewTestERC20(t.contractAddress, t.client)
		if err != nil {
			log.Printf("Failed to create contract instance: %v", err)
			return
		}
	}

	currentBalance, err := t.contract.BalanceOf(&bind.CallOpts{}, t.contractAddress)
	if err != nil {
		log.Printf("Failed to get contract balance: %v", err)
		return
	}

	totalSupply, err := t.contract.TotalSupply(&bind.CallOpts{})
	if err != nil {
		log.Printf("Failed to get total supply: %v", err)
		return
	}

	t.totalTransferredOut = new(big.Int).Sub(totalSupply, currentBalance)
}
