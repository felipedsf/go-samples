package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	EXCHANGE_SERVICE_URL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
)

var svc ExchangeService

func main() {
	fmt.Println("Server is running!")

	svc = ExchangeService{
		db: GetDatabase(),
	}

	http.HandleFunc("/cotacao", GetExchange)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetExchange(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	now := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, EXCHANGE_SERVICE_URL, nil)
	if err != nil {
		log.Printf("error creating request: %s\n", err.Error())
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error executing request: %s\n", err.Error())
		return
	}
	defer res.Body.Close()

	select {
	case <-time.After(200 * time.Millisecond):
		log.Printf("Exchange service called successfully: %s\n", time.Since(now))
	case <-ctx.Done():
		log.Println("timeout on call exchange service")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(ctx.Err().Error()))
		return
	}

	var exchange Exchange
	err = json.NewDecoder(res.Body).Decode(&exchange)
	if err != nil {
		log.Printf("error on decode json %s\n", err.Error())
		return
	}

	svc.InsertExchange(exchange)
	resp, err := json.Marshal(ExchangeResult{
		Bid: exchange.Usdbrl.Bid,
	})
	if err != nil {
		log.Printf("error to marshal json %s\n", err.Error())
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		log.Fatal(err)
		return
	}
}
