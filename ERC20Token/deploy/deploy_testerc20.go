package deploy

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ymytheresa/erc20-token-tracker/ERC20Token/connection"
	"github.com/ymytheresa/erc20-token-tracker/ERC20Token/contractsgo"
	"github.com/ymytheresa/erc20-token-tracker/ERC20Token/interact"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func deployTestERC20Contract() string {
	auth, client, fromAddress, nonce, gasPrice, _ := connection.GetNextTransaction()

	fmt.Println("Deploying TestERC20 contract...")

	auth.From = fromAddress
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(30000000)
	auth.GasPrice = gasPrice

	address, tx, _, err := contractsgo.DeployTestERC20(auth, client)
	if err != nil {
		fmt.Println("here")
		log.Fatal(err)
	}

	_, err = bind.WaitDeployed(context.Background(), client, tx)
	if err != nil {
		fmt.Println("there")
		log.Fatal(err)
	}

	fmt.Println("The contract is deployed at address: ", address)
	fmt.Printf("Transaction hash: 0x%x\n\n", tx.Hash())

	// Write address to file
	rootDir := filepath.Join("..") // Adjust this based on your project structure
	filePath := filepath.Join(rootDir, "contract_address.txt")
	err = os.WriteFile(filePath, []byte(address.String()), 0644)
	if err != nil {
		log.Fatalf("Failed to write contract address: %v", err)
	}

	fmt.Printf("Contract address saved to: %s\n", filePath)

	return address.String()
}

func RunTestERC20Contract() {
	testERC20ContractAddress := deployTestERC20Contract()
	toAddress := connection.GenerateNewWallet()
	interact.TransferTokens(testERC20ContractAddress, toAddress, 10)
}
