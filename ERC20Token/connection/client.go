package connection

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func GoDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func Connection() (*ethclient.Client, *big.Int, common.Address, *ecdsa.PrivateKey) {
	//poor quality code that the Connection restricted the transaction to the contract owner's address only
	client, err := ethclient.Dial(GoDotEnvVariable("GANACHE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("You are now connected to local Ganache!")

	privateKeyFromEnv := GoDotEnvVariable("DEPLOYER_PRIVATE_KEY")

	privateKey, err := crypto.HexToECDSA(privateKeyFromEnv)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("From address:", fromAddress.Hex())
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("The ETH balance of the account is:", balance)
	fmt.Println("------------------------------------------------------------------------")

	return client, chainID, fromAddress, privateKey //return contract owner's address as fromAddress
}

func GetClientForContractTx() *ethclient.Client {
	client, err := ethclient.Dial(GoDotEnvVariable("GANACHE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

// GetNextTransaction returns the next transaction in the pending transaction queue
func GetNextTransaction() (*bind.TransactOpts, *ethclient.Client, common.Address, uint64, *big.Int, error) {
	client, chainID, fromAddress, privateKey := Connection()

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// sign the transaction
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, client, fromAddress, nonce, gasPrice, err
	}

	return auth, client, fromAddress, nonce, gasPrice, nil
}

func GenerateNewWallet() common.Address {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return address
}
