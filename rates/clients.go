package rates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-redis/redis"
)

const KEY string = "rates"

// Database connection data definition
type databaseConnectionData struct {
	host     string
	port     string
	password string
}

// Interface for database operations
type databaseInterface interface {
	getConnectionData() *databaseConnectionData
	cacheRates(*exchangeRateResponsePayload) error
	fetchCachedRates() ([]byte, error)
}

// Endpoint (url + querystring) data
type endpointData struct {
	url         string
	queryString string
}

// Interface for web client operations
type webClientInterface interface {
	getEndpointData() *endpointData
	fetchLatestRates() ([]byte, error)
}

// Redis implementation of databaseInterface
type redisDatabase struct{}

// Exchange rates implementation of webClientInterface
type exchangeRatesAPI struct{}

// Returns database connection data
func (db *redisDatabase) getConnectionData() *databaseConnectionData {
	return &databaseConnectionData{
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		password: os.Getenv("DB_PASSWORD"),
	}
}

// Returns a new database client
func (db *redisDatabase) getConnection() (*redis.Client, error) {
	var connectionString bytes.Buffer
	config := db.getConnectionData()

	connectionString.WriteString(config.host)
	connectionString.WriteString(":")
	connectionString.WriteString(config.port)

	client := redis.NewClient(&redis.Options{
		Addr:     connectionString.String(),
		Password: config.password,
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return client, nil
}

// Returns cached rates payload
func (db *redisDatabase) fetchCachedRates() ([]byte, error) {
	client, err := db.getConnection()
	if err != nil {
		return nil, databaseError("could not connect to database")
	}
	defer client.Close()

	logger.Info("querying cache for rates")
	rates, err := client.Get(KEY).Bytes()
	if err != nil {
		logger.Error(err)
		return nil, databaseError(fmt.Sprintf("key not found: %s", KEY))
	}

	logger.Info("retrieved rates from cache")
	return rates, nil
}

// Caches rates payload
func (db *redisDatabase) cacheRates(payload *exchangeRateResponsePayload) error {
	client, err := db.getConnection()
	if err != nil {
		return databaseError("could not connect to database")
	}
	defer client.Close()

	data, _ := json.Marshal(payload)
	logger.Info(fmt.Sprintf("setting cache for %s", KEY))
	client.Set(KEY, data, 0)
	return nil
}

// Returns external rates api endpoint configuration
func (api *exchangeRatesAPI) getEndpointData() *endpointData {
	return &endpointData{
		url:         os.Getenv("EXCHANGE_RATES_URL"),
		queryString: os.Getenv("EXCHANGE_RATES_QUERYSTRING"),
	}
}

// Returns latest rates from exchange rates api
func (api *exchangeRatesAPI) fetchLatestRates() ([]byte, error) {
	config := api.getEndpointData()

	logger.Trace("requesting API for latest rates")
	request, err := http.Get(config.url + config.queryString)
	if err != nil {
		logger.Error(err)
		return nil, apiError("API request failed")
	}
	defer request.Body.Close()

	if request.StatusCode != http.StatusOK {
		return nil, apiError(fmt.Sprintf("HTTP Status code - %d", request.StatusCode))
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		logger.Error(err)
		return nil, apiError("could not parse body")
	}

	logger.Info("retrieved latest rates from api")
	return body, nil
}

// Returns a new redis database client according to databaseInterface
func newRedisDatabaseClient() databaseInterface {
	return &redisDatabase{}
}

// Returns a new exchange rates api client according to webClientInterface
func newExchangeRatesAPIClient() webClientInterface {
	return &exchangeRatesAPI{}
}
