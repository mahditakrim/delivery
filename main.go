package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

var orderQueue = make(chan deliveryOrderRequest, 10000)

var (
	orderID      = 0
	orderIDMutex = sync.Mutex{}
)

func orderIDGenerator() int {

	orderIDMutex.Lock()
	defer orderIDMutex.Unlock()
	orderID++
	return orderID
}

type deliveryOrderRequest struct {
	ID                int    `json:"id"`
	UserAddress       string `json:"user_address"`
	Origin            string `json:"origin"`
	Destination       string `json:"destination"`
	DeliveryTimestamp int    `json:"delivery_timestamp"`
	Status            string `json:"status"`
}

func main() {

	go startLogger()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handleQueue(ctx)
		close(done)
	}()

	go func() {
		http.HandleFunc("/core/order", handleOrder)
		http.HandleFunc("/delivery/webhook", handleDeliveryWebhook)
		fmt.Println("Server is running on localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	cancel()
	<-done
}
