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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, EXCHANGE_SERVICE_URL, nil)
	if err != nil {
		log.Printf("error creating request: %s\n", err.Error())
		return
	}
	defer req.Body.Close()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error executing request: %s\n", err.Error())
		return
	}

	select {
	case <-time.After(200 * time.Millisecond):
		log.Println("Exchange service called successfully")
	case <-ctx.Done():
		log.Println("timeout on call exchange service")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var exchange Exchange
	err = json.NewDecoder(res.Body).Decode(&exchange)
	if err != nil {
		log.Printf("error on decode json %s\n", err.Error())
		return
	}

	svc.InsertExchange(exchange)
	_, err = w.Write([]byte(fmt.Sprintf("{\"Bid\": \"%s\"}", exchange.Usdbrl.Bid)))
	if err != nil {
		log.Fatal(err)
		return
	}
}
