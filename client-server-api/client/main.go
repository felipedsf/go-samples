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
)

type Exchange struct {
	Bid string `json:"bid"`
}

func main() {
	fmt.Println("Client is running!")
	ctx, cncl := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cncl()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer req.Body.Close()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	select {
	case <-time.After(300 * time.Millisecond):
		log.Println("Request processed with success")
	case <-ctx.Done():
		log.Fatal("Cancelled by client")
		return
	}

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
}
