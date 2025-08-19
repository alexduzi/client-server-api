package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/alexduzi/client-server-api/model"
	"github.com/alexduzi/client-server-api/server/database"
)

const (
	URL_EXCHANGE_RATE    string        = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	TIMEOUT_EXCHANGE_API time.Duration = time.Millisecond * 200
	TIMEOUT_EXCHANGE_DB  time.Duration = time.Millisecond * 10
)

func main() {
	CreateServer()
}

func CreateServer() {
	database.InitializeDb()
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", ExchangeRateFunc)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func ExchangeRateFunc(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctxApi, cancelApi := context.WithTimeout(ctx, TIMEOUT_EXCHANGE_API)
	defer cancelApi()

	req, err := http.NewRequestWithContext(ctxApi, "GET", URL_EXCHANGE_RATE, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errMsg := fmt.Errorf("não foi possível obter cotação atual do dolar: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg.Error()))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	ctxDb, cancelDb := context.WithTimeout(ctx, TIMEOUT_EXCHANGE_DB)
	defer cancelDb()

	database.InsertExchange(ctxDb, URL_EXCHANGE_RATE, string(body))

	var exchange model.Exchange
	json.Unmarshal(body, &exchange)

	exchangeResp := model.ExchangeResponse{Bid: exchange.Usdbrl.Bid}
	exchangeRespJson, _ := json.Marshal(exchangeResp)

	w.WriteHeader(http.StatusOK)
	w.Write(exchangeRespJson)
}
