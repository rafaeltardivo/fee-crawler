package rates

import (
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
}

// Rates response.
type RatesResponse struct {
	Payload *ExchangeRatesResponsePayload
	Err     error
}

// Converts values to RatesResponse.
func toRatesResponse(ratesPayload *ExchangeRatesResponsePayload, err error) RatesResponse {
	return RatesResponse{
		Payload: ratesPayload,
		Err:     err,
	}
}
