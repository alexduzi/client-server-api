package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/alexduzi/client-server-api/model"
)

const (
	URL_EXCHANGE_RATE       string        = "http://localhost:8080/cotacao"
	TIMEOUT_EXCHANGE_CLIENT time.Duration = time.Millisecond * 300
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT_EXCHANGE_CLIENT)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", URL_EXCHANGE_RATE, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	saveContentToFile(body)
}

func saveContentToFile(body []byte) {
	var exchangeResp model.ExchangeResponse
	json.Unmarshal(body, &exchangeResp)

	file, err := os.OpenFile("./cotacao.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write([]byte(fmt.Sprintf("Dólar: %s", exchangeResp.Bid)))
	if err != nil {
		panic(err)
	}
}
