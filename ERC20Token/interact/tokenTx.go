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

func TransferTokens(contractAddress string, toAddress common.Address, value int64) (string, error) {
	_, client, fromAddress, nonce, gasPrice, _ := connection.GetNextTransaction() //fromAddress is contract owner's address

	fmt.Println("Transferring TestERC20 tokens...")
	txHash, err := transferTokensWithGasEstimate(client, fromAddress, toAddress, nonce, gasPrice, value, contractAddress)
	if err != nil {
		return "", err
	}
	return txHash.Hex(), nil
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

func transferTokensWithGasEstimate(client *ethclient.Client, fromAddress common.Address, toAddress common.Address, nonce uint64, gasPrice *big.Int, value int64, contractAddress string) (common.Hash, error) {
	gasLimit, err := estimateGasForTransfer(client, fromAddress, toAddress, contractAddress, value)
	if err != nil {
		return common.Hash{}, err
	}
	fmt.Println("Estimated gas:", gasLimit)

	auth, _, _, _, _, err := connection.GetNextTransaction()
	if err != nil {
		return common.Hash{}, err
	}

	auth.From = fromAddress
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice

	testERC20 := GetTestERC20Contract(client, contractAddress) //contractOwner consent, contract address
	contractAddressObj := common.HexToAddress(contractAddress)

	fmt.Println("\nBefore Transfer:")
	printAddressDetails(client, testERC20, "Contract", contractAddressObj)
	printAddressDetails(client, testERC20, "Sender", fromAddress)
	printAddressDetails(client, testERC20, "Receiver", toAddress)

	tx, err := testERC20.Transfer(auth, toAddress, big.NewInt(value)) //with contract owner consent and contract address, and toAddress, we now transfer tokens
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to transfer tokens: %v", err)
	}

	_, err = bind.WaitMined(context.Background(), client, tx) //wait for the transaction to be mined
	if err != nil {
		return common.Hash{}, err
	}
	fmt.Printf("\nTransaction hash: 0x%x\n", tx.Hash())

	fmt.Println("\nAfter Transfer:")
	printAddressDetails(client, testERC20, "Contract", contractAddressObj)
	printAddressDetails(client, testERC20, "Sender", fromAddress)
	printAddressDetails(client, testERC20, "Receiver", toAddress)

	return tx.Hash(), nil
}

func estimateGasForTransfer(client *ethclient.Client, fromAddress common.Address, toAddress common.Address, contractAddress string, value int64) (uint64, error) {
	store := abi.MustParseMethod("transfer(address,uint256)") //calling transfer function inside goeth abi for ERC20 contract

	abiData, err := store.EncodeArgs(toAddress, big.NewInt(value))
	if err != nil {
		return 0, err
	}

	toContractAddress := common.HexToAddress(contractAddress)

	callMsg := ethereum.CallMsg{
		From:     fromAddress,
		To:       &toContractAddress,	//this line always takes contract address
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

func printAddressDetails(client *ethclient.Client, testERC20 *contractsgo.TestERC20, label string, address common.Address) {
	ethBalance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Printf("Failed to get ETH balance for %s: %v", label, err)
		return
	}

	tokenBalance, err := testERC20.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Printf("Failed to get token balance for %s: %v", label, err)
		return
	}

	fmt.Printf("%s Address: %s\n", label, address.Hex())
	fmt.Printf("%s ETH Balance: %s wei\n", label, ethBalance.String())
	fmt.Printf("%s Token Balance: %s\n\n", label, tokenBalance.String())
}
