package rates

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestNormalizeFeeToBRL(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	ret := normalizeFeeToBRL("1,00")

	g.Expect(ret).Should(gomega.Equal("1.00"), "Return should have the same value and . as decimal separator")
}

func TestSanitizeFeeNotFound(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	ret, err := sanitizeFee("Sometext")

	g.Expect(ret).To(gomega.Equal(""), "Return should be empty")
	g.Expect(err).Should(gomega.HaveOccurred(), "An error should have occurred")
	g.Expect(err).Should(gomega.MatchError(parseError("fee not found")), "Error should be a parseError")
}

func TestSanitizeFee(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	ret, err := sanitizeFee("Some 1,00 text")

	g.Expect(ret).Should(gomega.Equal("1,00"), "Return should be equal to 1,00")
	g.Expect(err).ShouldNot(gomega.HaveOccurred(), "An error should not have occurred")
}

func TestToCurrencyValue(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	ret := toCurrencyValue("6.00", 0.3333333333)

	g.Expect(ret).Should(gomega.Equal("1.99"), "Return should be equal to 1.99")
}

func TestToRateData(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	exchangeRateResponseMock := exchangeRateResponsePayload{
		Rates: exchangeRateCurrencyPayload{
			EUR: 0.1687023416,
			USD: 0.1875295229,
		},
		Base: "BRL",
		Date: "2020-06-01",
	}

	rateDataMock := rateData{
		BRL:  "1.00",
		USD:  toCurrencyValue("1.00", exchangeRateResponseMock.Rates.USD),
		EUR:  toCurrencyValue("1.00", exchangeRateResponseMock.Rates.EUR),
		Date: exchangeRateResponseMock.Date,
	}

	convertedRateData, err := toRateData("1,00", &exchangeRateResponseMock)

	g.Expect(convertedRateData).Should(gomega.Equal(&rateDataMock), "Return should be equal to rateDataMock")
	g.Expect(err).ShouldNot(gomega.HaveOccurred(), "An error should not have occurred")
}
