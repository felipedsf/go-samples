package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	URL       = "http://localhost:8080/cotacao"
	FILE_NAME = "cotacao.txt"
	TIMEOUT   = 300 * time.Millisecond
)

type Exchange struct {
	Bid string `json:"bid"`
}

func main() {
	fmt.Println("Client is running!")
	now := time.Now()

	ctx, cncl := context.WithTimeout(context.Background(), TIMEOUT)
	defer cncl()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var exchange Exchange
	err = json.NewDecoder(res.Body).Decode(&exchange)
	if err != nil {
		log.Fatal("error to unmarshal json ", err)
	}

	file, err := os.Create(FILE_NAME)
	if err != nil {
		log.Fatal("error to create file ", err)
	}
	defer file.Close()
	_, err = file.Write([]byte(fmt.Sprintf("Dólar: %s", exchange.Bid)))
	if err != nil {
		log.Fatal("error to write file", err)
	}
	log.Printf("executed successfully in %s - bid: %s\n", time.Since(now), exchange.Bid)
}
