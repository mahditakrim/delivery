package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func handleOrder(w http.ResponseWriter, r *http.Request) {

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

	req.ID = orderIDGenerator()
	orderQueue <- req
	logStatus("pending")

	resp, _ := json.Marshal(req)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}
