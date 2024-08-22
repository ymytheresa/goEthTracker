package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ymytheresa/erc20-token-tracker/ERC20Token/interact"
)

var randomAddresses []common.Address

func RandomTransaction(interval time.Duration, done chan bool) {
	numberOfRandomAddresses := 10
	randomAddresses := generateRandomAddresses(numberOfRandomAddresses)
	defer printAllAddresses()

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		contractAddr, err := getContractAddress()
		if err != nil {
			log.Println("Error getting contract address:", err)
			done <- true
			return
		}

		for {
			select {
			case <-ticker.C:
				transact(contractAddr, randomAddresses)
			case <-done:
				return
			}
		}
	}()
}

func transact(contractAddr string, randomAddresses []common.Address) {
	interact.TransferTokens(contractAddr, randomAddresses[rand.Intn(len(randomAddresses))], int64(rand.Intn(100)))
}

func getContractAddress() (string, error) {
	filePath := "./contract_address.txt"
	address, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(address), nil
}

func generateRandomAddresses(n int) []common.Address {
	addresses := make([]common.Address, n)
	for i := 0; i < n; i++ {
		addresses[i] = common.HexToAddress(fmt.Sprintf("0x%x", rand.Uint64()))
	}
	randomAddresses = addresses
	return addresses
}

func printAllAddresses() {
	for _, address := range randomAddresses {
		fmt.Println(address.Hex())
	}
}
