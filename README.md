# Ethereum-based dApp Development with Go using go-ethereum SDK

## Examples deploying smart contracts on Neon EVM Devnet using Go

This directory contains all the files necessary to deploy the following smart contracts on Neon EVM Devnet-

1. Simple storage contract.
2. ERC20 token contract.

For more details, please refer to these documentations https://goethereumbook.org/en/ and https://geth.ethereum.org/docs/developers/dapp-developer/native-bindings.

## Prerequisites

1. The latest Go version.
2. Solidity compiler version <= 0.8.25 (Neon EVM supports solidity <= 0.8.26 but Homebrew only supports 0.8.25 for now).

### Solc installation

Please check this link to install the required solc version - https://docs.soliditylang.org/en/latest/installing-solidity.html.

### Go installation

> [!IMPORTANT]
> If your machine already has Go installed, then please proceed to the **Cloning repository** step.

1. Download the latest Go version from https://go.dev/doc/install.

2. Create a directory `GoProjects` for your Go project development on your machine. Please run `pwd` inside `GoProjects` directory to get the full path. This will be used in **Step 4** to set the $GOPATH env variable.

3. Create 3 directories `/GoProjects/src`, `/GoProjects/pkg` and `/GoProjects/bin`.

4. Set up the `$GOPATH` env variable on your machine -

- Run `nano ~/.zshrc` on Mac machines.
- Run `nano ~/.bash_profile` on Linux machines.
- Paste the following lines -

```sh
export GOPATH=<PATH_TO_YOUR_GO_PROJECTS_DIRECTORY>
export PATH=$GOPATH/bin:$PATH
```

5. Save your `~/.bash_profile` or `~/.zshrc` file.

- Run `source ~/.zshrc` on Mac machines.
- Run `source ~/.bash_profile` on Linux machines.

6. Run `echo $GOPATH` to check if the GOPATH is set correctly in the machine.

> [!IMPORTANT]
> Neon EVM doesn't support the latest JSON-RPC specifications. Therefore, Neon EVM only supports `go-ethereum` versions **<=1.12.2**.

## Cloning repository

1. Navigate to the `src` directory of your GOPATH.

```sh
cd $GOPATH/src
```

2. Run command

```sh
git clone https://github.com/neonlabsorg/goEthTracker.git
cd goEthTracker
```

**Note:** All the following commands should be executed inside `goEthTracker` folder.

## Install the required libraries inside the `goEthTracker` folder

```sh
npm install
```

```sh
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```

```sh
go mod tidy
```

## Set up .env file

Rename `.env.example` to `.env` and place your private key inside it.

## Interact with the **Storage** smart contract

### Generate the go bindings

1. Run the following commands to generate the smart contract ABI and bytecode.

```sh
solc --abi ./contracts/Storage.sol -o build
```

```sh
solc --bin ./contracts/Storage.sol -o build
```

2. Run the following command to generate the smart contract go binding inside `contractsgo` folder.

```sh
 abigen --abi ./build/Storage.abi --pkg contractsgo --type Storage --out ./contractsgo/Storage.go --bin ./build/Storage.bin
```

### Run the smart contracts functions

> [!IMPORTANT]
> To run only the Storage contract, please comment out `deploy.RunTestERC20Contract()` in `main.go`.

Run the following command to deploy the Storage contract, store a value in the deployed smart contract and reading the value from the deployed smart contract.

```sh
go run main.go
```

After successfully running this step you should get console output similar to:

```sh
You are now connected to Neon EVM Devnet
The NEON balance of the account is: 310387553758242748088682
------------------------------------------------------------------------
Deploying Storage contract...
The contract is deployed at address:  0x6b6Ba862e2bBc0C1305DF681d45f16a1D6F57baf
Transaction hash: 0xf84667ce0bd5d2da4dfcf81fe9c72bdc81c207a41a3c9baa4c43e9ebb6ae1b6e

You are now connected to Neon EVM Devnet
The NEON balance of the account is: 310383249542814769793482
------------------------------------------------------------------------
Storing value in the Storage contract...
Estimated gas: 25000
Transaction hash: 0x24e5af83df1e9f1536d684c08e903d1285f1f5e484df43d4616c925bb25ec9a9

You are now connected to Neon EVM Devnet
The NEON balance of the account is: 310383247282862115123482
------------------------------------------------------------------------
Reading value from the Storage contract...
Returned value: 45
```

## Interact with the **TestERC20** smart contract

### Generate the go bindings

1. Run the following commands to generate the smart contract ABI and bytecode.

```sh
solc --abi ./contracts/TestERC20.sol -o build
```

```sh
solc --bin ./contracts/TestERC20.sol -o build
```

2. Run the following command to generate the smart contract go binding inside `contractsgo` folder.

```sh
 abigen --abi ./build/TestERC20.abi --pkg contractsgo --type TestERC20 --out ./contractsgo/TestERC20.go --bin ./build/TestERC20.bin
```

### Run the smart contracts functions

> [!IMPORTANT]
> To run only the TestERC20 contract, please comment out `deploy.RunStorageContract()` in `main.go`.

Run the following command to deploy the TestERC20 contract and transfer some TestERC20 tokens from the deployer address to a randomly created address

```sh
go run main.go
```

After successfully running this step you should get console output similar to:

```sh
You are now connected to Neon EVM Devnet
The NEON balance of the account is: 310383247282862115123482
------------------------------------------------------------------------
Deploying TestERC20 contract...
The contract is deployed at address:  0x7BeE8180c4f35744C9cC811e540252ECcD8AcEb4
Transaction hash: 0xf8af65bcb8187bcdcc8c7a5a7106f242c941d6506201497f31f46099d891bcc6

You are now connected to Neon EVM Devnet
The NEON balance of the account is: 310373551028315738437922
------------------------------------------------------------------------
Transferring TestERC20 tokens...
Estimated gas: 1422000
Sender balance before transfer: 1000000000000000000000
Receiver balance before transfer: 0
Transaction hash: 0x8d2ff2a94f836b25e3ae9cc2f9b95ca73e3b3c1e4a6bf7725890eddd915029ab

Sender balance after transfer: 999999999999999999990
Receiver balance after transfer: 10
```
