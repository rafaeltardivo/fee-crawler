package rates

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// Interface for rate parsers
type parserInterface interface {
	ParseRates([]byte) (interface{}, error)
}

// Exchange rates implementation of parser interface
type exchangeRatesParser struct{}

type exchangeRateCurrencyPayload struct {
	EUR float64 `json:"EUR"`
	USD float64 `json:"USD"`
}

// Structured exchange rates payload
type exchangeRateResponsePayload struct {
	Rates exchangeRateCurrencyPayload `json:"rates"`
	Base  string                      `json:"base"`
	Date  string                      `json:"date"`
}

// Parses the received payload and returns a type-compatible interface
func (p *exchangeRatesParser) ParseRates(data []byte) (interface{}, error) {
	var payload exchangeRateResponsePayload

	err := json.Unmarshal(data, &payload)
	if err != nil {
		return nil, parseError(err.Error())
	}

	// Input range validation (OWASP recommendation)
	if payload.Rates.EUR <= 0 || payload.Rates.USD <= 0 {
		return nil, parseError(fmt.Sprintf("Invalid rate data"))
	}

	// Input format validation (OWASP recommendation)
	re := regexp.MustCompile("((19|20)\\d\\d)-(0?[1-9]|[12][0-9]|3[01])-(0?[1-9]|1[012])")
	if re.MatchString(payload.Date) == false {
		return nil, parseError(fmt.Sprintf("Invalid date format"))
	}

	// Input domain validation (OWASP recommendation)
	if payload.Base != "BRL" {
		return nil, parseError(fmt.Sprintf("Invalid base rate"))
	}

	return &payload, nil
}

// Creates and returns a new exchange rates parser according to parserInterface
func newExchangeRatesParser() parserInterface {
	return &exchangeRatesParser{}
}
