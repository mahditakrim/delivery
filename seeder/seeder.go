package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type deliveryOrderRequest struct {
	UserAddress       string `json:"user_address"`
	Origin            string `json:"origin"`
	Destination       string `json:"destination"`
	DeliveryTimestamp int    `json:"delivery_timestamp"`
}

func main() {

	for {
		time.Sleep(time.Millisecond * 200)

		go func() {
			timestamp := rand.Intn(2*60*60) + int(time.Now().Unix())
			body, err := json.Marshal(deliveryOrderRequest{
				DeliveryTimestamp: timestamp,
			})
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return
			}

			r, _ := http.NewRequest("POST", "http://localhost:8080/core/order", bytes.NewBuffer(body))
			r.Header.Add("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(r)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
			defer func() {
				_ = resp.Body.Close()
			}()

			body, err = io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response:", err)
				return
			}

			fmt.Println("Response Body:", string(body))
		}()
	}
}
