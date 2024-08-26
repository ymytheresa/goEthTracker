package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/flock"
	"github.com/ymytheresa/erc20-token-tracker/ERC20Token/interact"
)

var randomAddresses []common.Address
var fileLock *flock.Flock
var mutex sync.Mutex

func init() {
	fileLock = flock.New("hash.txt")
}

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
				txHash := transact(contractAddr, randomAddresses) //doesnt mean its transferring from contract address, its transferring from contract owner's address
				if err := writeTransactionHash(txHash); err != nil {
					log.Printf("Error writing transaction hash: %v", err)
				}
			case <-done:
				return
			}
		}
	}()
}

func transact(contractAddr string, randomAddresses []common.Address) string {
	recipient := randomAddresses[rand.Intn(len(randomAddresses))]
	txHash, err := interact.TransferTokens(contractAddr, recipient, int64(rand.Intn(100))) //transferring tokens from contract owner's address to random address actually. but contract address is needed
	if err != nil {
		log.Printf("Error in transaction: %v", err)
		return "" // Return an empty string in case of error
	}
	return txHash
}

func writeTransactionHash(txHash string) error {
	if txHash == "" {
		return nil
	}

	mutex.Lock()
	defer mutex.Unlock()

	locked, err := fileLock.TryLock()
	if err != nil {
		return fmt.Errorf("error acquiring file lock: %v", err)
	}
	if !locked {
		return fmt.Errorf("could not acquire file lock")
	}
	defer fileLock.Unlock()

	file, err := os.OpenFile("hash.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(txHash + "\n"); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
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
