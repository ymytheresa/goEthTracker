# Crypto Airdrop Tracker

This CLI tool allows you to track crypto airdrops by token owner using a local Ethereum network (Ganache).

## Prerequisites

Before you can use this tool, you need to install the following:

1. Go (Golang)
2. Ganache

### Installing Go

1. Visit the official Go download page: https://golang.org/dl/
2. Download the installer for your operating system.
3. Follow the installation instructions for your OS.
4. Verify the installation by opening a terminal and running:
   ```
   go version
   ```

### Installing Ganache

1. Visit the Ganache website: https://www.trufflesuite.com/ganache
2. Download the appropriate version for your operating system.
3. Install Ganache following the instructions for your OS.
4. Launch Ganache to ensure it's working correctly.

## Usage

This tool consists of several Go scripts that simulate and track token transactions on a local Ethereum network.

### Deploying the ERC20 Token Contract

Run the following command to deploy the TestERC20 contract:
(put the owner address in owner_address.txt and copy the .env.example to .env in all subdirectories where main.go lies)

```
go run ./ERC20Token
```

This will:
- Connect to local Ganache
- Deploy the TestERC20 contract
- Save the contract address to `contract_address.txt`

### Simulating Airdrops

To simulate airdrops, run:

```
go run ./TokenTransaction
```

This script will:
- Generate 10 random Ethereum addresses
- Transfer 1-100 random TestERC20 tokens to these addresses 

### Tracking Airdrops

To track the airdrops, run:

```
go run ./TokenTracker
```

This will:
- Connect to the local Ganache network
- Monitor token transfers
- Display interval and total sums of tokens transferred to each address

## Sample Output

The tool provides detailed output at each step. Here's an example of what you might see:

```
You are now connected to local Ganache!
From address: 0xa652010de06D0C0E6d589289C11bC1D7914191d9
The ETH balance of the account is: 999998104250000000000
------------------------------------------------------------------------
Interval Sums:
0x0000000000000000000000007f471C251808AFf0: 3
0x0000000000000000000000002EbbAD7e4A6eaf40: 23
0x000000000000000000000000D86694EF9A06518c: 88

Total Sums:
0x000000000000000000000000D86694EF9A06518c: 88
0x0000000000000000000000007f471C251808AFf0: 3
0x0000000000000000000000002EbbAD7e4A6eaf40: 23
...
```

## Note

This tool is designed for educational and testing purposes on a local Ethereum network. Always exercise caution when dealing with real cryptocurrencies and tokens.
Test cases not added yet. will add

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[Insert your chosen license here]
