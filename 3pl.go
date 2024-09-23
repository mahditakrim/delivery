package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	STATUS_SEARCHING = "searching"
	STATUS_FOUND     = "found"
	STATUS_NOT_FOUND = "not_found"
	STATUS_DELIVERED = "delivered"
)

var notFoundHistory = make(map[int]int)
var notFoundHistoryMutex = sync.Mutex{}

func call3plAPI(req deliveryOrderRequest) error {

	time.Sleep(time.Second)
	if rand.Intn(100) < 5 {
		return errors.New("3pl API error")
	}

	go func() {
		callDeliveryWebhook(req.ID, STATUS_SEARCHING)
		time.Sleep(time.Second)
		notFoundHistoryMutex.Lock()
		defer notFoundHistoryMutex.Unlock()
		notFoundCount := notFoundHistory[req.ID]
		if notFoundCount == 3 {
			callDeliveryWebhook(req.ID, STATUS_FOUND)
			delete(notFoundHistory, req.ID)
			statusMutex.Lock()
			notFound -= 3
			statusMutex.Unlock()
		} else {
			if rand.Intn(100) < 10 {
				notFoundHistory[req.ID]++
				callDeliveryWebhook(req.ID, STATUS_NOT_FOUND)
				return
			} else {
				callDeliveryWebhook(req.ID, STATUS_FOUND)
				statusMutex.Lock()
				notFound -= notFoundCount
				statusMutex.Unlock()
				delete(notFoundHistory, req.ID)
			}
		}
		time.Sleep(time.Second)
		callDeliveryWebhook(req.ID, STATUS_DELIVERED)
	}()

	return nil
}

func callDeliveryWebhook(orderID int, status string) {

	body, _ := json.Marshal(deliveryOrderRequest{
		ID:     orderID,
		Status: status,
	})
	r, _ := http.NewRequest("POST", "http://localhost:8080/delivery/webhook", bytes.NewBuffer(body))
	r.Header.Add("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(r)
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
}
