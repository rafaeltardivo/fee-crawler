package api

import (
	"github.com/rafaeltardivo/fee-crawler/crawler"
	"github.com/rafaeltardivo/fee-crawler/rates"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
}

// Truncates converted value as a decimal and returns a .xx "precision" string
func toCurrencyValue(fee string, rate float64) string {
	feeAmount, _ := decimal.NewFromString(fee)
	rateAmount := decimal.NewFromFloat(rate)

	return feeAmount.Mul(rateAmount).Truncate(2).String()
}

// Parses and converts values to FeeStruct
func toAPIResponse(ratesResponse rates.RatesResponse, crawlerResponse crawler.CrawlerResponse) (*feeStruct, error) {

	if ratesResponse.Err != nil {
		return nil, ratesResponse.Err
	}

	if crawlerResponse.Err != nil {
		return nil, crawlerResponse.Err
	}

	return &feeStruct{
		RatesDate:   ratesResponse.Payload.Date,
		Description: crawlerResponse.Description,
		BRL:         crawlerResponse.Amount,
		USD:         toCurrencyValue(crawlerResponse.Amount, ratesResponse.Payload.Rates.USD),
		EUR:         toCurrencyValue(crawlerResponse.Amount, ratesResponse.Payload.Rates.EUR),
	}, nil
}
