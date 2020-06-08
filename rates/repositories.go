package rates

import "fmt"

// Interface for rates repository
type repositoryInterface interface {
	fetchRates(databaseInterface, webClientInterface) (*ExchangeRatesResponsePayload, error)
}

type exchangeRatesRepository struct{}

// Returns rates payload
// Action 1: Look for today rates on database
// Action 2: Request today rates to exchange rates API
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

// Returns a new exchange rates repository according to repositoryInterface
func newExchangeRatesRepository() repositoryInterface {
	return &exchangeRatesRepository{}
}

// Rates command responsible for returning rates (from cache or api)
func GetRatesData() (*ExchangeRatesResponsePayload, error) {

	db := newRedisDatabaseClient()
	api := newExchangeRatesAPIClient()
	repository := newExchangeRatesRepository()

	logger.Info("starting rates command")
	return repository.fetchRates(db, api)
}
