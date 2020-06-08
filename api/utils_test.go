package api

import (
	"errors"
	"testing"

	"github.com/onsi/gomega"
	"github.com/rafaeltardivo/fee-crawler/crawler"
	"github.com/rafaeltardivo/fee-crawler/rates"
)

func TestToAPIResponseRatesResponseError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	ratesResponse := rates.RatesResponse{
		Payload: &rates.ExchangeRatesResponsePayload{
			Rates: rates.ExchangeRatesCurrencyPayload{
				EUR: 0.1687023416,
				USD: 0.1875295229,
			},
			Base: "BRL",
			Date: "2020-06-01",
		},
		Err: errors.New("some mocked rates error"),
	}
	crawlerResponse := crawler.CrawlerResponse{
		Amount:      "7.00",
		Description: "Some description",
		Err:         nil,
	}

	apiResponse, err := toAPIResponse(ratesResponse, crawlerResponse)

	g.Expect(apiResponse).To(gomega.BeNil(), "Return should be nil")
	g.Expect(err).To(gomega.HaveOccurred(), "An error should have occurred")
	g.Expect(err).To(gomega.MatchError(errors.New("some mocked rates error")), "Error should be match mocked one")
}

func TestToAPIResponseCrawlerResponseError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	ratesResponse := rates.RatesResponse{
		Payload: &rates.ExchangeRatesResponsePayload{
			Rates: rates.ExchangeRatesCurrencyPayload{
				EUR: 0.1687023416,
				USD: 0.1875295229,
			},
			Base: "BRL",
			Date: "2020-06-01",
		},
		Err: nil,
	}
	crawlerResponse := crawler.CrawlerResponse{
		Amount:      "7.00",
		Description: "Some description",
		Err:         errors.New("some mocked crawler error"),
	}

	apiResponse, err := toAPIResponse(ratesResponse, crawlerResponse)

	g.Expect(apiResponse).To(gomega.BeNil(), "Return should be nil")
	g.Expect(err).To(gomega.HaveOccurred(), "An error should have occurred")
	g.Expect(err).To(gomega.MatchError(errors.New("some mocked crawler error")), "Error should be match mocked one")
}

func TestToAPIResponse(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	ratesResponse := rates.RatesResponse{
		Payload: &rates.ExchangeRatesResponsePayload{
			Rates: rates.ExchangeRatesCurrencyPayload{
				EUR: 0.1687023416,
				USD: 0.1875295229,
			},
			Base: "BRL",
			Date: "2020-06-01",
		},
		Err: nil,
	}
	crawlerResponse := crawler.CrawlerResponse{
		Amount:      "7.00",
		Description: "Some description",
		Err:         nil,
	}

	apiResponse, err := toAPIResponse(ratesResponse, crawlerResponse)

	g.Expect(err).ToNot(gomega.HaveOccurred(), "An error should not have occurred")
	g.Expect(apiResponse.RatesDate).To(gomega.Equal(ratesResponse.Payload.Date), "Date should match ratesResponse value")
	g.Expect(apiResponse.Description).To(gomega.Equal(crawlerResponse.Description), "Description should match crawlerResponse value")
	g.Expect(apiResponse.BRL).To(gomega.Equal(crawlerResponse.Amount), "BRL should match ratesPayload value")

	currencyUSDValue := toCurrencyValue(apiResponse.BRL, ratesResponse.Payload.Rates.USD)
	g.Expect(apiResponse.USD).To(gomega.Equal(currencyUSDValue), "USD should match converted ratesPayload currency value")
	currencyEURValue := toCurrencyValue(apiResponse.BRL, ratesResponse.Payload.Rates.EUR)
	g.Expect(apiResponse.EUR).To(gomega.Equal(currencyEURValue), "EUR should match converted ratesPayload currency value")
}

func TestToCurrencyValue(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	value := toCurrencyValue("5.00", 0.7500000000)

	g.Expect(value).To(gomega.Equal("3.75"), "USD should match converted ratesPayload currency value")
}
