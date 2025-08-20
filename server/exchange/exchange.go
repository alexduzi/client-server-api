package exchange

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	URL_EXCHANGE_RATE    string        = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	TIMEOUT_EXCHANGE_API time.Duration = time.Millisecond * 200 // 200ms
)

func GetExchangeRate(ctx context.Context) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT_EXCHANGE_API)
	defer cancel()

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

	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
		return nil, ctx.Err()
	default:
		log.Println("api exchange call success")
	}

	return body, nil
}
