package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ymytheresa/erc20-token-tracker/ERC20Token/contractsgo"
	"github.com/ymytheresa/erc20-token-tracker/ERC20Token/interact"
)

// TransferEvent represents a simplified transfer event
type TransferEvent struct {
	From   string
	To     string
	Value  string
	TxHash string
}

// TokenTracker handles the tracking of ERC20 token events
type TokenTracker struct {
	contract            *contractsgo.TestERC20
	client              *ethclient.Client
	contractAddress     common.Address
	ownerAddress        common.Address
	transfers           []TransferEvent
	totalTransferredOut *big.Int
}

func setTokenTracker() TokenTracker {
	contractAddress, err := GetContractAddress()
	if err != nil {
		log.Fatal(err)
	}

	ownerAddress, err := GetOwnerAddress()
	if err != nil {
		log.Fatal(err)
	}

	client := interact.GetClient()

	// Verify connection to the correct network
	networkID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("Failed to get network ID:", err)
	}
	fmt.Printf("Connected to network with ID: %s\n", networkID.String())

	return TokenTracker{
		client:              client,
		contractAddress:     contractAddress,
		ownerAddress:        ownerAddress,
		transfers:           []TransferEvent{},
		totalTransferredOut: big.NewInt(0),
	}
}

// StartTracking begins listening for Transfer events
func (t *TokenTracker) StartTracking() error {
	txHash := common.HexToHash("0xcb6d95b941d3213eeaaaf3d8b78cec8b4d71400f9712bc649a1e64d59124ab85")
	tx, isPending, err := t.client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %v", err)
	}

	if isPending {
		fmt.Println("Transaction is still pending")
	} else {
		receipt, err := t.client.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			return fmt.Errorf("failed to get transaction receipt: %v", err)
		}

		fmt.Printf("Transaction details:\n")
		from, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
		if err != nil {
			return fmt.Errorf("failed to get transaction sender: %v", err)
		}
		fmt.Printf("From: %s\n", from.Hex())
		fmt.Printf("To: %s\n", tx.To().Hex())
		fmt.Printf("Value: %s wei\n", tx.Value().String())
		fmt.Printf("Gas Price: %s wei\n", tx.GasPrice().String())
		fmt.Printf("Gas Limit: %d\n", tx.Gas())
		fmt.Printf("Nonce: %d\n", tx.Nonce())
		fmt.Printf("Block Number: %d\n", receipt.BlockNumber)
		fmt.Printf("Gas Used: %d\n", receipt.GasUsed)

		for _, log := range receipt.Logs {
			event, err := t.parseTransferEvent(*log)
			if err != nil {
				fmt.Printf("Failed to parse event: %v\n", err)
				continue
			}
			fmt.Printf("Transfer Event: From %s, To %s, Value %s\n", event.From, event.To, event.Value)
		}
	}

	return nil
}

// parseTransferEvent parses a log entry into a TransferEvent
func (t *TokenTracker) parseTransferEvent(vLog types.Log) (TransferEvent, error) {
	event := TransferEvent{
		From:   vLog.Topics[1].Hex(),
		To:     vLog.Topics[2].Hex(),
		Value:  new(big.Int).SetBytes(vLog.Data).String(),
		TxHash: vLog.TxHash.Hex(),
	}
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

func GetOwnerAddress() (common.Address, error) {
	filePath := "./owner_address.txt"
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
		fmt.Printf("Total amount transferred from owner: %s\n", t.totalTransferredOut.String())
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

	ownerBalance, err := t.contract.BalanceOf(&bind.CallOpts{}, t.ownerAddress)
	if err != nil {
		log.Printf("Failed to get owner balance: %v", err)
		return
	}

	totalSupply, err := t.contract.TotalSupply(&bind.CallOpts{})
	if err != nil {
		log.Printf("Failed to get total supply: %v", err)
		return
	}

	t.totalTransferredOut = new(big.Int).Sub(totalSupply, ownerBalance)

	fmt.Println("\nContract and Owner Details:")
	fmt.Printf("Contract Address: %s\n", t.contractAddress.Hex())
	fmt.Printf("Owner Address: %s\n", t.ownerAddress.Hex())
	fmt.Printf("Owner Balance: %s\n", ownerBalance.String())
	fmt.Printf("Total Supply: %s\n", totalSupply.String())
	fmt.Printf("Total Transferred Out: %s\n", t.totalTransferredOut.String())

	fmt.Println("\nRecent Transfers from Owner:")
	transferCount := 0
	for i := len(t.transfers) - 1; i >= 0 && transferCount < 5; i-- {
		event := t.transfers[i]
		if event.From == t.ownerAddress.Hex() {
			fmt.Printf("To: %s, Amount: %s, TxHash: %s\n", event.To, event.Value, event.TxHash)
			transferCount++
		}
	}

	if transferCount == 0 {
		fmt.Println("No recent transfers from owner")
	}
}
