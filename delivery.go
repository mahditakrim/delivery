package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func handleDeliveryWebhook(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var req deliveryOrderRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Error unmarshal JSON", http.StatusBadRequest)
		return
	}

	go func() {
		logStatus(req.Status)
		if req.Status == STATUS_NOT_FOUND {
			logStatus("pending")
		RETRY:
			err := call3plAPI(req)
			if err != nil {
				goto RETRY
			}
		}
	}()

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func handleQueue(ctx context.Context) {

	wp := newWorkerPool(10)

	for {
		select {
		case <-ctx.Done():
			wp.shutdown()
			return
		case order := <-orderQueue:
			wp.launch(func() {
				t := time.Unix(int64(order.DeliveryTimestamp), 0)
				callTime := t.Add(-time.Hour)
				cmp := time.Now().Compare(callTime)
				if cmp == 0 || cmp == 1 {
				RETRY:
					err := call3plAPI(order)
					if err != nil {
						goto RETRY
					}
					return
				}
				schedule(time.Until(callTime), order)
			})
		}
	}
}

func schedule(d time.Duration, order deliveryOrderRequest) {

	time.AfterFunc(d, func() {
	RETRY:
		err := call3plAPI(order)
		if err != nil {
			goto RETRY
		}
	})
}
