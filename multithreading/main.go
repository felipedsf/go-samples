package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	BRASIL_API_URL = "https://brasilapi.com.br/api/cep/v1/#CEP#"
	VIA_CEP_URL    = "http://viacep.com.br/ws/#CEP#/json/"
)

var (
	cep string
	wg  *sync.WaitGroup
)

func main() {
	flag.StringVar(&cep, "cep", "", "CEP para ser consultado")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	brasilApiUrl := strings.ReplaceAll(BRASIL_API_URL, "#CEP#", cep)
	viaCepUrl := strings.ReplaceAll(VIA_CEP_URL, "#CEP#", cep)

	brasilApiReq, err := http.NewRequestWithContext(ctx, http.MethodGet, brasilApiUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	viaCepReq, err := http.NewRequest(http.MethodGet, viaCepUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	brasilCh := make(chan map[string]interface{})
	viaCh := make(chan map[string]interface{})
	wg := sync.WaitGroup{}
	wg.Add(1)

	go DoCallDec(viaCh, viaCepReq)
	go DoCallDec(brasilCh, brasilApiReq)

	select {
	case msg := <-brasilCh:
		log.Printf("BrasilApi: %s", msg)
		wg.Done()
	case msg := <-viaCh:
		log.Printf("ViaCep: %s", msg)
		wg.Done()
	case <-ctx.Done():
		log.Fatal("timeout")
	}
	wg.Wait()
	log.Println("Executed...")
}

func DoCallDec(c chan map[string]interface{}, req *http.Request) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var j map[string]any
	err = json.NewDecoder(resp.Body).Decode(&j)
	if err != nil {
		log.Fatal(err)
	}
	c <- j
}
