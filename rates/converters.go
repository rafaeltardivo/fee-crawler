package rates

// Definition for rate data
type rateData struct {
	BRL  string `json:"BRL"`
	USD  string `json:"USD"`
	EUR  string `json:"EUR"`
	Date string `json:"date"`
}

// Interface for rate converters
type converterInterface interface {
	ConvertRates(string, clientInterface, parserInterface) (*rateData, error)
}

// Exchange rates implementation of converter interface
type exchangeRatesConverter struct{}

// Convert exchange rates data response to structured data
func (c *exchangeRatesConverter) ConvertRates(
	fee string, client clientInterface, parser parserInterface) (*rateData, error) {

	data, err := client.GetRates()
	if err != nil {
		return nil, err
	}

	parsedData, err := parser.ParseRates(data)
	if err != nil {
		return nil, err
	}

	exchangeRateData, _ := parsedData.(*exchangeRateResponsePayload)
	convertedData, err := toRateData(fee, exchangeRateData)
	if err != nil {
		return nil, err
	}

	return convertedData, nil
}

// Creates and returns a new exchange rates converter according to converterInterface
func NewExchangeRatesConverter() converterInterface {
	return &exchangeRatesConverter{}
}
