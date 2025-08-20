package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/alexduzi/client-server-api/model"
	"github.com/alexduzi/client-server-api/server/database"

	"github.com/alexduzi/client-server-api/server/exchange"
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
	body, err := exchange.GetExchangeRate(context.Background())
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = database.InsertExchange(context.Background(), exchange.URL_EXCHANGE_RATE, string(body))
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var exchange model.Exchange
	json.Unmarshal(body, &exchange)

	exchangeResp := model.ExchangeResponse{Bid: exchange.Usdbrl.Bid}
	exchangeRespJson, _ := json.Marshal(exchangeResp)

	w.WriteHeader(http.StatusOK)
	w.Write(exchangeRespJson)
}
