package rates

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/go-redis/redis"
	"github.com/onsi/gomega"
)

func TestNewExchangeRatesRepository(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	repository := newExchangeRatesRepository()
	_, isOne := repository.(repositoryInterface)

	g.Expect(isOne).To(gomega.BeTrue(), "Repository is-one repositoryInterface")
}

type mockedDatabaseCachedRates struct{}

func (m *mockedDatabaseCachedRates) getConnectionData() *databaseConnectionData {
	return &databaseConnectionData{}
}
func (m *mockedDatabaseCachedRates) getConnection() (*redis.Client, error) {
	return nil, nil
}
func (m *mockedDatabaseCachedRates) cacheRates(payload *ExchangeRatesResponsePayload) error {
	return nil
}
func (m *mockedDatabaseCachedRates) fetchCachedRates() ([]byte, error) {
	responseMock := ExchangeRatesResponsePayload{
		Rates: ExchangeRatesCurrencyPayload{
			EUR: 0.1687023416,
			USD: 0.1875295229,
		},
		Base: "BRL",
		Date: "2020-06-01",
	}
	rawBody := new(bytes.Buffer)
	json.NewEncoder(rawBody).Encode(responseMock)
	return rawBody.Bytes(), nil
}

func TestRepositoryCachedRates(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	expectedResponse := &ExchangeRatesResponsePayload{
		Rates: ExchangeRatesCurrencyPayload{
			EUR: 0.1687023416,
			USD: 0.1875295229,
		},
		Base: "BRL",
		Date: "2020-06-01",
	}

	db := &mockedDatabaseCachedRates{}
	api := newExchangeRatesAPIClient()
	repository := newExchangeRatesRepository()

	ret, err := repository.fetchRates(db, api)
	g.Expect(err).ToNot(gomega.HaveOccurred(), "An error should not have occurred")
	g.Expect(ret).To(gomega.Equal(expectedResponse), "Return should be equal to expectedResponse")
}

type mockedDatabaseCachedError struct{}

func (m *mockedDatabaseCachedError) getConnectionData() *databaseConnectionData {
	return &databaseConnectionData{}
}
func (m *mockedDatabaseCachedError) getConnection() (*redis.Client, error) {
	return nil, nil
}
func (m *mockedDatabaseCachedError) cacheRates(payload *ExchangeRatesResponsePayload) error {
	return nil
}
func (m *mockedDatabaseCachedError) fetchCachedRates() ([]byte, error) {
	return nil, databaseError("Some error")
}

type mockedAPILatestRates struct{}

func (m *mockedAPILatestRates) getEndpointData() *endpointData {
	return nil
}
func (m *mockedAPILatestRates) fetchLatestRates() ([]byte, error) {
	responseMock := ExchangeRatesResponsePayload{
		Rates: ExchangeRatesCurrencyPayload{
			EUR: 0.1687023416,
			USD: 0.1875295229,
		},
		Base: "BRL",
		Date: "2020-06-01",
	}
	rawBody := new(bytes.Buffer)
	json.NewEncoder(rawBody).Encode(responseMock)
	return rawBody.Bytes(), nil
}

func TestRepositoryLatestRates(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	expectedResponse := &ExchangeRatesResponsePayload{
		Rates: ExchangeRatesCurrencyPayload{
			EUR: 0.1687023416,
			USD: 0.1875295229,
		},
		Base: "BRL",
		Date: "2020-06-01",
	}

	db := &mockedDatabaseCachedError{}
	api := &mockedAPILatestRates{}
	repository := newExchangeRatesRepository()

	ret, err := repository.fetchRates(db, api)
	g.Expect(err).ToNot(gomega.HaveOccurred(), "An error should not have occurred")
	g.Expect(ret).To(gomega.Equal(expectedResponse), "Return should be equal to expectedResponse")
}

type mockedDatabaseCachedInvalidPayload struct{}

func (m *mockedDatabaseCachedInvalidPayload) getConnectionData() *databaseConnectionData {
	return &databaseConnectionData{}
}
func (m *mockedDatabaseCachedInvalidPayload) getConnection() (*redis.Client, error) {
	return nil, nil
}
func (m *mockedDatabaseCachedInvalidPayload) cacheRates(payload *ExchangeRatesResponsePayload) error {
	return nil
}
func (m *mockedDatabaseCachedInvalidPayload) fetchCachedRates() ([]byte, error) {
	return []byte("some invalid payload"), nil
}

func TestRepositoryInvalidPayloadParseError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	db := &mockedDatabaseCachedInvalidPayload{}
	api := newExchangeRatesAPIClient()
	repository := newExchangeRatesRepository()

	ret, err := repository.fetchRates(db, api)
	g.Expect(ret).To(gomega.BeNil(), "Return should be nil")
	g.Expect(err).To(gomega.HaveOccurred(), "An error should have occurred")
	g.Expect(err).To(gomega.MatchError(parseError("could not parse payload")), "Error should be a parseError")
}

type mockedAPILatestRatesError struct{}

func (m *mockedAPILatestRatesError) getEndpointData() *endpointData {
	return nil
}
func (m *mockedAPILatestRatesError) fetchLatestRates() ([]byte, error) {
	return nil, apiError("some mocked api error")
}
func TestRepositoryInvalidCacheAndLatestRates(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	db := &mockedDatabaseCachedError{}
	api := &mockedAPILatestRatesError{}
	repository := newExchangeRatesRepository()

	ret, err := repository.fetchRates(db, api)
	g.Expect(ret).To(gomega.BeNil(), "Return should be nil")
	g.Expect(err).To(gomega.HaveOccurred(), "An error should have occurred")
	g.Expect(err).To(gomega.MatchError(apiError("some mocked api error")), "Error should be a parseError")
}
