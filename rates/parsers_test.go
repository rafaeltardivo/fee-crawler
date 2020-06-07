package rates

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/onsi/gomega"
)

func TestnewExchangeRatesParserType(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	parser := newExchangeRatesParser()
	_, isOne := parser.(parserInterface)

	g.Expect(isOne).To(gomega.BeTrue(), "parser is-one parserInterface")
}

func TestParserRatesRatesError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	exchangeRateResponseMock := exchangeRateResponsePayload{
		Rates: exchangeRateCurrencyPayload{
			EUR: 0.00,
			USD: 0.00,
		},
		Base: "BRL",
		Date: "2020-06-01",
	}
	rawBody := new(bytes.Buffer)
	json.NewEncoder(rawBody).Encode(exchangeRateResponseMock)

	parser := newExchangeRatesParser()
	ret, err := parser.ParseRates(rawBody.Bytes())

	g.Expect(ret).Should(gomega.BeNil(), "Return should be nil")
	g.Expect(err).Should(gomega.HaveOccurred(), "An error should have occurred")
	g.Expect(err).Should(gomega.MatchError(parseError("invalid rate data")), "Error should be a parseError")
}

func TestParserRatesDateError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	exchangeRateResponseMock := exchangeRateResponsePayload{
		Rates: exchangeRateCurrencyPayload{
			EUR: 0.1687023416,
			USD: 0.1875295229,
		},
		Base: "BRL",
		Date: "2020-06",
	}
	rawBody := new(bytes.Buffer)
	json.NewEncoder(rawBody).Encode(exchangeRateResponseMock)

	parser := newExchangeRatesParser()
	ret, err := parser.ParseRates(rawBody.Bytes())

	g.Expect(ret).Should(gomega.BeNil(), "Return should be nil")
	g.Expect(err).Should(gomega.HaveOccurred(), "An error should have occurred")
	g.Expect(err).Should(gomega.MatchError(parseError("invalid date format")), "Error should be a parseError")
}

func TestParserRatesBaseError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	exchangeRateResponseMock := exchangeRateResponsePayload{
		Rates: exchangeRateCurrencyPayload{
			EUR: 0.1687023416,
			USD: 0.1875295229,
		},
		Base: "USD",
		Date: "2020-06-01",
	}
	rawBody := new(bytes.Buffer)
	json.NewEncoder(rawBody).Encode(exchangeRateResponseMock)

	parser := newExchangeRatesParser()
	ret, err := parser.ParseRates(rawBody.Bytes())

	g.Expect(ret).Should(gomega.BeNil(), "Return should be nil")
	g.Expect(err).Should(gomega.HaveOccurred(), "An error should have occurred")
	g.Expect(err).Should(gomega.MatchError(parseError("invalid base rate")), "Error should be a parseError")
}

func TestParserRates(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	exchangeRateResponseMock := exchangeRateResponsePayload{
		Rates: exchangeRateCurrencyPayload{
			EUR: 0.1687023416,
			USD: 0.1875295229,
		},
		Base: "BRL",
		Date: "2020-06-01",
	}
	rawBody := new(bytes.Buffer)
	json.NewEncoder(rawBody).Encode(exchangeRateResponseMock)

	parser := newExchangeRatesParser()
	ret, err := parser.ParseRates(rawBody.Bytes())
	exchangeRateData, _ := ret.(*exchangeRateResponsePayload)

	g.Expect(err).ShouldNot(gomega.HaveOccurred(), "An error should not have occurred")
	g.Expect(ret).Should(gomega.Equal(exchangeRateData), "Return should be equal to exchangeRateData")
}
