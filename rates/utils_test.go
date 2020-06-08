package rates

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestToRatesResponse(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	ratesPayload := &ExchangeRatesResponsePayload{
		Rates: ExchangeRatesCurrencyPayload{
			EUR: 0.1687023416,
			USD: 0.1875295229,
		},
		Base: "BRL",
		Date: "2020-06-01",
	}
	ratesResponse := toRatesResponse(ratesPayload, nil)

	g.Expect(ratesResponse.Payload.Rates.EUR).To(gomega.Equal(ratesPayload.Rates.EUR), "EUR should match ratesPayload value")
	g.Expect(ratesResponse.Payload.Rates.USD).To(gomega.Equal(ratesPayload.Rates.USD), "USD should match ratesPayload value")
	g.Expect(ratesResponse.Payload.Base).To(gomega.Equal(ratesPayload.Base), "USD should match ratesPayload value")
	g.Expect(ratesResponse.Payload.Date).To(gomega.Equal(ratesPayload.Date), "Date should match ratesPayload value")
}
