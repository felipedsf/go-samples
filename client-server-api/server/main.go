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

func main() {
	fmt.Println("Server is running!")

	http.HandleFunc("/usd", GetExchange)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetExchange(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, EXCHANGE_SERVICE_URL, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
		return
	}
	_, err = w.Write([]byte(fmt.Sprintf("{\"Bid\": \"%s\"}", exchange.Usdbrl.Bid)))
	if err != nil {
		log.Fatal(err)
		return
	}
}

type Exchange struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}
