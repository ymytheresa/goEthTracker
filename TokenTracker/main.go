package main

import "log"

func main() {
	// tracker := setTokenTracker()
	// go tracker.StartTracking() // Start the existing event tracking

	// // Monitor contract transfers every 5 minutes
	// go tracker.MonitorContractTransfers(5 * time.Second)

	// // Keep the main goroutine running
	// select {}
	tracker := setTokenTracker()
	err := tracker.StartTracking()
	if err != nil {
		log.Fatal(err)
	}
}
