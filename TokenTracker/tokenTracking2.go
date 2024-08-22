package main

import (
	"bufio"
	"context"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ymytheresa/erc20-token-tracker/ERC20Token/interact"
)

var (
	hashFilePath = "hash.txt"
	mutex        sync.Mutex
	intervalSums = make(map[common.Address]*big.Int)
	totalSums    = make(map[common.Address]*big.Int)
	mapMutex     sync.Mutex
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
			printMaps()
			resetIntervalSums()
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
		events, err := getEventsForHash(txHash)
		if err != nil {
			fmt.Printf("Error getting events for hash %s: %v\n", txHash, err)
			continue
		}
		// printEvents(events)
		updateMaps(events)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading hash file: %v\n", err)
		return err
	}

	// Clear the file after processing
	os.Truncate(hashFilePath, 0)
	return nil
}

func getEventsForHash(txHash string) ([]TransferEvent, error) {
	hash := common.HexToHash(txHash)
	client := interact.GetClient()
	receipt, err := client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return nil, fmt.Errorf("error getting transaction receipt: %v", err)
	}

	var events []TransferEvent
	for _, log := range receipt.Logs {
		if len(log.Topics) == 3 && log.Topics[0] == common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef") {
			from := common.HexToAddress(log.Topics[1].Hex())
			to := common.HexToAddress(log.Topics[2].Hex())
			value := new(big.Int).SetBytes(log.Data)

			event := TransferEvent{
				From:   from,
				To:     to,
				Value:  value,
				TxHash: hash,
			}
			events = append(events, event)
		}
	}

	return events, nil
}

func printEvents(events []TransferEvent) {
	for _, event := range events {
		fmt.Printf("Transaction Hash: %s\n", event.TxHash.Hex())
		fmt.Printf("From: %s\n", event.From.Hex())
		fmt.Printf("To: %s\n", event.To.Hex())
		fmt.Printf("Value: %s\n", event.Value.String())
		fmt.Println("------------------------")
	}
}

func updateMaps(events []TransferEvent) {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	for _, event := range events {
		// Update intervalSums
		if _, exists := intervalSums[event.To]; !exists {
			intervalSums[event.To] = big.NewInt(0)
		}
		intervalSums[event.To].Add(intervalSums[event.To], event.Value)

		// Update totalSums
		if _, exists := totalSums[event.To]; !exists {
			totalSums[event.To] = big.NewInt(0)
		}
		totalSums[event.To].Add(totalSums[event.To], event.Value)
	}
}

func printMaps() {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	fmt.Println("Interval Sums:")
	for addr, sum := range intervalSums {
		fmt.Printf("%s: %s\n", addr.Hex(), sum.String())
	}

	fmt.Println("\nTotal Sums:")
	for addr, sum := range totalSums {
		fmt.Printf("%s: %s\n", addr.Hex(), sum.String())
	}
	fmt.Println()
}

func resetIntervalSums() {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	intervalSums = make(map[common.Address]*big.Int)
}
