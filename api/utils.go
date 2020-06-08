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

// Truncates converted value as a decimal and returns a .xx "precision" string.
func toCurrencyValue(fee string, rate float64) string {
	feeAmount, _ := decimal.NewFromString(fee)
	rateAmount := decimal.NewFromFloat(rate)

	return feeAmount.Mul(rateAmount).Truncate(2).String()
}

// Parses and converts values to FeePayload.
func toAPIResponse(ratesResponse rates.RatesResponse, crawlerResponse crawler.CrawlerResponse) (*feePayload, error) {

	if ratesResponse.Err != nil {
		return nil, ratesResponse.Err
	}

	if crawlerResponse.Err != nil {
		return nil, crawlerResponse.Err
	}

	return &feePayload{
		RatesDate:   ratesResponse.Payload.Date,
		Description: crawlerResponse.Description,
		BRL:         crawlerResponse.Amount,
		USD:         toCurrencyValue(crawlerResponse.Amount, ratesResponse.Payload.Rates.USD),
		EUR:         toCurrencyValue(crawlerResponse.Amount, ratesResponse.Payload.Rates.EUR),
	}, nil
}

// Validates domain input
func validateDomain(domain string) bool {
	validDomains := []string{"http://www.smartmei.com.br", "https://www.smartmei.com.br"}
	for _, item := range validDomains {
		if domain == item {
			return true
		}
	}
	return false
}
