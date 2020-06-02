package rates

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Interface for rate clients
type clientInterface interface {
	GetRates() ([]byte, error)
}

// ExchangeRates client implementation of clientInterface
type exchangeRatesClient struct {
	url     string
	filters string
}

// Requests currency data from REST api and returns byte array containig the raw response
func (c *exchangeRatesClient) GetRates() ([]byte, error) {
	var endpoint bytes.Buffer
	endpoint.WriteString(c.url)
	endpoint.WriteString(c.filters)
	request, err := http.Get(endpoint.String())
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

	if request.StatusCode != http.StatusOK {
		return nil, clientError(fmt.Sprintf("HTTP Status code: %d", request.StatusCode))
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Creates and returns a new exchange rates client according to clientInterface
func NewExchangeRatesClient() clientInterface {
	return &exchangeRatesClient{
		url:     "https://api.exchangeratesapi.io/latest",
		filters: "?base=BRL&symbols=USD,EUR",
	}
}
