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
	URL_EXCHANGE_RATE       string        = "http://localhost:8080/contacao"
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

	var exchange model.Exchange
	json.Unmarshal(body, &exchange)

	file, err := os.OpenFile("./cotacao.txt", os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(fmt.Sprintf("Dólar: %f", exchange.Usdbrl.Bid))
}
