package interact

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/ymytheresa/erc20-token-tracker/ERC20Token/connection"
	"github.com/ymytheresa/erc20-token-tracker/ERC20Token/contractsgo"

	"github.com/defiweb/go-eth/abi"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	testERC20Contract *contractsgo.TestERC20
	testERC20Mu       sync.Mutex
)

func TransferTokens(contractAddress string, toAddress common.Address, value int64) {
	_, client, fromAddress, nonce, gasPrice, _ := connection.GetNextTransaction()

	fmt.Println("Transferring TestERC20 tokens...")
	transferTokensWithGasEstimate(client, fromAddress, toAddress, nonce, gasPrice, value, contractAddress)
}

func GetClient() *ethclient.Client {
	_, client, _, _, _, _ := connection.GetNextTransaction()
	return client
}

func GetTestERC20Contract(client *ethclient.Client, contractAddress string) *contractsgo.TestERC20 {
	testERC20, err := contractsgo.NewTestERC20(common.HexToAddress(contractAddress), client)
	if err != nil {
		log.Fatalf("Failed to instantiate TestERC20 contract: %v", err)
	}
	return testERC20
}

func transferTokensWithGasEstimate(client *ethclient.Client, fromAddress common.Address, toAddress common.Address, nonce uint64, gasPrice *big.Int, value int64, contractAddress string) {
	gasLimit, err := estimateGasForTransfer(client, fromAddress, toAddress, contractAddress, value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Estimated gas:", gasLimit)

	auth, _, _, _, _, err := connection.GetNextTransaction()
	if err != nil {
		log.Fatal(err)
	}

	auth.From = fromAddress
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice

	testERC20 := GetTestERC20Contract(client, contractAddress)

	fmt.Println("Sender address:", fromAddress.String())
	fmt.Println("Receiver address:", toAddress.String())
	fmt.Println("Sender balance before transfer:", GetBalance(testERC20, fromAddress))
	fmt.Println("Receiver balance before transfer:", GetBalance(testERC20, toAddress))

	tx, err := testERC20.Transfer(auth, toAddress, big.NewInt(value))
	if err != nil {
		log.Fatalf("Failed to transfer tokens: %v", err)
	}

	_, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction hash: 0x%x\n\n", tx.Hash())

	fmt.Println("Sender balance after transfer:", GetBalance(testERC20, fromAddress))
	fmt.Println("Receiver balance after transfer:", GetBalance(testERC20, toAddress))
}

func estimateGasForTransfer(client *ethclient.Client, fromAddress common.Address, toAddress common.Address, contractAddress string, value int64) (uint64, error) {
	store := abi.MustParseMethod("transfer(address,uint256)")
	// Encode method arguments.
	abiData, err := store.EncodeArgs(toAddress, big.NewInt(value))
	if err != nil {
		return 0, err
	}

	toContractAddress := common.HexToAddress(contractAddress)

	callMsg := ethereum.CallMsg{
		From:     fromAddress,
		To:       &toContractAddress,
		GasPrice: nil,
		Value:    big.NewInt(0),
		Data:     abiData,
	}

	return client.EstimateGas(context.Background(), callMsg)
}

func GetBalance(testERC20 *contractsgo.TestERC20, address common.Address) *big.Int {
	balance, err := testERC20.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}
	return balance
}
