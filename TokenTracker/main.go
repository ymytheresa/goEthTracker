package main

import (
	"fmt"
	"log"
)

func main() {
	tracker, err := NewTokenTracker()
	if err != nil {
		log.Fatalf("Failed to create token tracker: %v", err)
	}

	fmt.Printf("Starting to track events for contract: %s\n", tracker.contractAddress.Hex())

	err = tracker.StartTracking()
	if err != nil {
		log.Fatalf("Failed to start tracking: %v", err)
	}

	// Keep the program running
	select {}
}
