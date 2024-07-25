package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/felipedsf/go-samples/client-server-api/server/db"
	"github.com/felipedsf/go-samples/client-server-api/server/service"
	"log"
	"net/http"
	"time"
)

const (
	EXCHANGE_SERVICE_URL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	TIMEOUT              = 200 * time.Millisecond
)

var svc service.ExchangeService

func main() {
	fmt.Println("Server is running!")

	svc = service.ExchangeService{
		Db: db.GetDatabase(),
	}

	http.HandleFunc("/cotacao", GetExchange)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetExchange(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

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

	var exchange service.Exchange
	err = json.NewDecoder(res.Body).Decode(&exchange)
	if err != nil {
		log.Printf("error on decode json %s\n", err.Error())
		return
	}

	svc.InsertExchange(exchange)
	resp, err := json.Marshal(service.ExchangeResult{
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
	log.Printf("executed successfully in %s - bid: %s\n", time.Since(now), exchange.Usdbrl.Bid)
}
