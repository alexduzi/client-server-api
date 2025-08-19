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
	TIMEOUT_EXCHANGE_API time.Duration = time.Millisecond * 200 // 200ms
	TIMEOUT_EXCHANGE_DB  time.Duration = time.Millisecond * 10  // 10 ms
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
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_EXCHANGE_API)
	defer cancel()

	body, err := getExchangeRate(ctx)
	if err != nil {
		errMsg := fmt.Errorf("não foi possível obter cotação atual do dolar: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg.Error()))
		return
	}

	select {
	case <-ctx.Done():
		errMsg := fmt.Errorf("timeout ao chamar api de cotação")
		log.Println(errMsg.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg.Error()))
		return
	default:
		log.Println("api de cotação chamada com sucesso")
	}

	ctx, cancel = context.WithTimeout(context.Background(), TIMEOUT_EXCHANGE_DB)
	defer cancel()

	err = database.InsertExchange(ctx, URL_EXCHANGE_RATE, string(body))
	if err != nil {
		errMsg := fmt.Errorf("timeout ao chamar api de cotação")
		log.Println(errMsg.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg.Error()))
	}

	select {
	case <-ctx.Done():
		errMsg := fmt.Errorf("timeout ao inserir cotação atual na base de dados")
		log.Println(errMsg.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg.Error()))
		return
	default:
		log.Println("cotação inserida com sucesso na base de dados")
	}

	var exchange model.Exchange
	json.Unmarshal(body, &exchange)

	exchangeResp := model.ExchangeResponse{Bid: exchange.Usdbrl.Bid}
	exchangeRespJson, _ := json.Marshal(exchangeResp)

	w.WriteHeader(http.StatusOK)
	w.Write(exchangeRespJson)
}

func getExchangeRate(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", URL_EXCHANGE_RATE, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
