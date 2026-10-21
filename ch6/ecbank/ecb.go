package ecbank

import (
	"fmt"
	"learning/ch6-moneyconverter/money"
	"net/http"
)

const (
	ErrCallingServer = ecbankError("error calling server")
)

// Client can call the bank to retrieve exchange rates.
type Client struct {
}

// FetchExchangeRate fetches today's ExchangeRate and returns it.
func (c Client) FetchExchangeRate(source, target money.Currency) (money.ExchangeRate, error) {
	const euroxrefURL = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

	resp, err := http.Get(euroxrefURL)
	if err != nil {
		return money.ExchangeRate{}, fmt.Errorf("%w: %s", ErrServerSide, err.Error())
	}
	defer resp.Body.Close()

	return money.ExchangeRate{}, nil
}
