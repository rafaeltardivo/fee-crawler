package rates

import (
	"fmt"
	"sync"
)

// Interface for rates repository.
type repositoryInterface interface {
	fetchRates(databaseInterface, webClientInterface) (*ExchangeRatesResponsePayload, error)
}

// Exchange rates implementation of repositoryInterface.
type exchangeRatesRepository struct{}

// Returns rates payload.
// Action 1: Look for today rates on database.
// Action 2: Request today rates to exchange rates API.
func (r *exchangeRatesRepository) fetchRates(db databaseInterface, api webClientInterface) (*ExchangeRatesResponsePayload, error) {
	updateCache := false

	raw, err := db.fetchCachedRates()
	if err != nil {
		logger.Warn("could not find cache, will now make an external request")

		updateCache = true
		raw, err = api.fetchLatestRates()
		if err != nil {
			return nil, err
		}
	}

	parser := newExchangeRatesParser()
	parsedData, err := parser.ParseRates(raw)
	if err != nil {
		return nil, err
	}
	rates, _ := parsedData.(*ExchangeRatesResponsePayload)

	if updateCache {
		db.cacheRates(rates)
	}

	logger.Info(fmt.Sprintf("retrieved rates: %+v", *rates))
	return rates, nil
}

// Returns a new exchange rates repository.
func newExchangeRatesRepository() repositoryInterface {
	return &exchangeRatesRepository{}
}

// Rates command responsible for managing rates requests.
// Passes the returned rates to channel.
func GetRatesData(done chan RatesResponse, wg *sync.WaitGroup) {
	defer close(done)
	defer wg.Done()

	db := newRedisDatabaseClient()
	api := newExchangeRatesAPIClient()
	repository := newExchangeRatesRepository()

	ratesPayload, err := repository.fetchRates(db, api)
	done <- toRatesResponse(ratesPayload, err)
}
