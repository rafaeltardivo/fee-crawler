package rates

import (
	"os"
	"testing"

	"github.com/onsi/gomega"
)

func TestNewRedisDatabaseClient(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	client := newRedisDatabaseClient()
	_, isOne := client.(databaseInterface)

	g.Expect(isOne).To(gomega.BeTrue(), "redis client is-one databaseInterface")
}

func TestNewExchangeRatesAPIClient(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	client := newExchangeRatesAPIClient()
	_, isOne := client.(webClientInterface)

	g.Expect(isOne).To(gomega.BeTrue(), "exchange rates client is-one webClientInterface")
}

func TestRedisDatabaseClientConnectionData(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	client := newRedisDatabaseClient()
	configurationData := client.getConnectionData()

	g.Expect(configurationData.host).To(gomega.Equal(os.Getenv("DB_HOST")), "configuration host is equal to env var DB_HOST")
	g.Expect(configurationData.port).To(gomega.Equal(os.Getenv("DB_PORT")), "configuration host is equal to env var DB_PORT")
	g.Expect(configurationData.password).To(gomega.Equal(os.Getenv("DB_PASSWORD")), "configuration host is equal to env var PASSWORD")
}

func TestExchangeRatesClientEndpointData(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	client := newExchangeRatesAPIClient()
	endpointData := client.getEndpointData()

	g.Expect(endpointData.url).To(gomega.Equal(os.Getenv("EXCHANGE_RATES_URL")), "url is equal to env var EXCHANGE_RATES_URL")
	g.Expect(endpointData.queryString).To(gomega.Equal(os.Getenv("EXCHANGE_RATES_QUERYSTRING")), "queryString is equal to env var EXCHANGE_RATES_QUERYSTRING")
}
