package main

import (
	"bufio"
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ymytheresa/erc20-token-tracker/ERC20Token/contractsgo"
	// Import your contract ABI package here
)

var (
	hashFilePath = "hash.txt"
	mutex        sync.Mutex
)

type TransferEvent struct {
	From   common.Address
	To     common.Address
	Value  *big.Int
	TxHash common.Hash
}

func startTicker() {
	fmt.Println("Starting ticker...")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		fmt.Println("Tick at", time.Now())
		select {
		case <-ticker.C:
			err := processTransactions()
			if err != nil {
				fmt.Printf("Error processing transactions: %v\n", err)
			}
		}
	}
}

func processTransactions() error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.Open(hashFilePath)
	if err != nil {
		fmt.Printf("Error opening hash file: %v\n", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txHash := scanner.Text()
		printEventForHash(txHash)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading hash file: %v\n", err)
		return err
	}

	// Clear the file after processing
	os.Truncate(hashFilePath, 0)
	return nil
}

func printEventForHash(txHash string) {
	client, err := ethclient.Dial("http://localhost:8545") // Connect to your local Ganache
	if err != nil {
		fmt.Printf("Error connecting to Ganache: %v\n", err)
		return
	}
	defer client.Close()

	receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	if err != nil {
		fmt.Printf("Error getting transaction receipt for %s: %v\n", txHash, err)
		return
	}

	contractABI, err := abi.JSON(strings.NewReader(contractsgo.TestERC20ABI))
	if err != nil {
		fmt.Printf("Error parsing contract ABI: %v\n", err)
		return
	}

	fmt.Printf("Events for transaction %s:\n", txHash)
	for _, log := range receipt.Logs {
		event, err := parseTransferEvent(contractABI, *log)
		if err != nil {
			fmt.Printf("Error parsing event: %v\n", err)
			continue
		}
		fmt.Printf("Transfer: From %s To %s Amount %s\n", event.From, event.To, event.Value)
	}
}

func parseTransferEvent(contractABI abi.ABI, log types.Log) (TransferEvent, error) {
	event := TransferEvent{}
	err := contractABI.UnpackIntoInterface(&event, "Transfer", log.Data)
	if err != nil {
		return event, err
	}
	event.From = common.HexToAddress(log.Topics[1].Hex())
	event.To = common.HexToAddress(log.Topics[2].Hex())
	event.TxHash = log.TxHash
	return event, nil
}
