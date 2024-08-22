package main

import "go-ethereum-tutorial/deploy"

func main() {
	// Storage contract
	deploy.RunStorageContract()

	// TestERC20 contract
	deploy.RunTestERC20Contract()
}