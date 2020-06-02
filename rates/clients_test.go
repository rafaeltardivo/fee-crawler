package rates

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestNewExchangeRatesClientType(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	client := NewExchangeRatesClient()
	_, sameType := client.(clientInterface)

	g.Expect(sameType).To(gomega.BeTrue(), "client is-one clientInterface")
}
