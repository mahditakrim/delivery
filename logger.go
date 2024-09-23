package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	pending   = 0
	searching = 0
	found     = 0
	notFound  = 0
	delivered = 0
)
var statusMutex = sync.Mutex{}

func logStatus(status string) {

	switch status {
	case STATUS_SEARCHING:
		logSearching()
	case STATUS_FOUND:
		logFound()
	case STATUS_DELIVERED:
		logDelivered()
	case STATUS_NOT_FOUND:
		logNotFound()
	default:
		logPending()
	}
}

func logPending() {

	statusMutex.Lock()
	defer statusMutex.Unlock()
	pending++
}

func logSearching() {

	statusMutex.Lock()
	defer statusMutex.Unlock()
	pending--
	searching++
}

func logFound() {

	statusMutex.Lock()
	defer statusMutex.Unlock()
	searching--
	found++
}

func logNotFound() {

	statusMutex.Lock()
	defer statusMutex.Unlock()
	searching--
	notFound++
}

func logDelivered() {

	statusMutex.Lock()
	defer statusMutex.Unlock()
	found--
	delivered++
}

func printStatus() {

	fmt.Print("\033[H\033[2J")
	fmt.Println("Pending:", pending)
	fmt.Println("Searching:", searching)
	fmt.Println("Not Found:", notFound)
	fmt.Println("Found:", found)
	fmt.Println("Delivered:", delivered)
}

func startLogger() {

	ticker := time.Tick(time.Second)
	for range ticker {
		printStatus()
	}
}
