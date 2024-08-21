package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/erc20-token-tracker/internal/config"
	"github.com/yourusername/erc20-token-tracker/internal/ethereum"
	"github.com/yourusername/erc20-token-tracker/internal/tracker"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize Ethereum client
	client, err := ethereum.NewClient(cfg.GanacheURL)
	if err != nil {
		return fmt.Errorf("failed to create Ethereum client: %w", err)
	}
	defer client.Close()

	// Deploy contract (this could be moved to a separate command)
	contractAddress, err := ethereum.DeployContract(client, cfg.DeployerPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to deploy contract: %w", err)
	}
	log.Printf("Contract deployed at address: %s", contractAddress.Hex())

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize and start the token tracker
	tokenTracker, err := tracker.NewTracker(client, contractAddress)
	if err != nil {
		return fmt.Errorf("failed to create token tracker: %w", err)
	}
	go tokenTracker.Start(ctx)

	// Wait for interrupt signal
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-shutdownSignal

	log.Println("Shutting down gracefully...")
	cancel()

	// Allow some time for graceful shutdown
	time.Sleep(5 * time.Second)
	log.Println("Shutdown complete")

	return nil
}
