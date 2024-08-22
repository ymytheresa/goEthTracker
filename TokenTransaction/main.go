package main

import "time"

func main() {
	done := make(chan bool)
	RandomTransaction(5*time.Second, done)
	<-done
}
