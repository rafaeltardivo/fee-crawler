package rates

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/onsi/gomega"
)

func TestNewExchangeRatesConverterType(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	converter := NewExchangeRatesConverter()
	_, sameType := converter.(converterInterface)

	g.Expect(sameType).To(gomega.BeTrue(), "converter is-one clientInterface")
}

type clientErrorMock struct{}

func (c *clientErrorMock) GetRates() ([]byte, error) {
	return nil, clientError("Mocked")
}

func TestConvertRatesClientError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	client := &clientErrorMock{}
	parser := NewExchangeRatesParser()
	converter := NewExchangeRatesConverter()

	ret, err := converter.ConvertRates("1.00", client, parser)

	g.Expect(ret).Should(gomega.BeNil(), "Return should be nil")
	g.Expect(err).Should(gomega.HaveOccurred(), "An error should have occurred")
	g.Expect(err).Should(gomega.MatchError(clientError("Mocked")), "Error should be a clientError")
}

type clientMock struct{}

func (c *clientMock) GetRates() ([]byte, error) {
	return []byte{}, nil
}

type parserErrorMock struct{}

func (c *parserErrorMock) ParseRates([]byte) (interface{}, error) {
	return nil, parseError("Mocked")
}

func TestConvertRatesParserError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	client := &clientMock{}
	parser := &parserErrorMock{}
	converter := NewExchangeRatesConverter()

	ret, err := converter.ConvertRates("1.00", client, parser)

	g.Expect(ret).Should(gomega.BeNil(), "Return should be nil")
	g.Expect(err).Should(gomega.HaveOccurred(), "An error should have occurred")
	g.Expect(err).Should(gomega.MatchError(parseError("Mocked")), "Error should be a parseError")
}

type parserMock struct{}

func (c *parserMock) ParseRates([]byte) (interface{}, error) {
	return rateData{}, nil
}

func TestConvertRatesConverterError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	client := &clientMock{}
	parser := &parserMock{}
	converter := NewExchangeRatesConverter()

	ret, err := converter.ConvertRates("Some Text", client, parser)

	g.Expect(ret).Should(gomega.BeNil(), "Return should be nil")
	g.Expect(err).Should(gomega.HaveOccurred(), "An error should have occurred")
	g.Expect(err).Should(gomega.MatchError(convertError("Fee not found")), "Error should be a convertError")
}

type ClientResponseMock struct{}

var exchangeRateResponseMock = exchangeRateResponsePayload{
	Rates: exchangeRateCurrencyPayload{
		EUR: 0.1687023416,
		USD: 0.1875295229,
	},
	Base: "BRL",
	Date: "2020-06-01",
}

func (c *ClientResponseMock) GetRates() ([]byte, error) {
	rawBody := new(bytes.Buffer)
	json.NewEncoder(rawBody).Encode(exchangeRateResponseMock)

	return rawBody.Bytes(), nil
}

func TestConvertRatesConverter(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	client := &ClientResponseMock{}
	parser := NewExchangeRatesParser()
	converter := NewExchangeRatesConverter()

	rateData, _ := toRateData("1,00", &exchangeRateResponseMock)

	ret, err := converter.ConvertRates("1,00", client, parser)
	g.Expect(ret).Should(gomega.Equal(rateData), "Return should be equal to rateData")
	g.Expect(err).ShouldNot(gomega.HaveOccurred(), "An error should not have occurred")
}
