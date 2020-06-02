package rates

import (
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
)

// Normalizes and returns the fee value using BRL floating point criteria
func normalizeFeeToBRL(fee string) string {
	return strings.Replace(fee, ",", ".", 1)
}

// Sanitizes and returns the raw crawled fee string
func sanitizeFee(fee string) (string, error) {
	re := regexp.MustCompile(`[0-9,]+[0-9]{2}`)
	rawValue := re.FindStringSubmatch(fee)

	if len(rawValue) <= 0 {
		return "", convertError("Fee not found")
	}

	return rawValue[0], nil
}

// Truncates converted value as a decimal and returns a .xx "precision" string
func toCurrencyValue(fee string, rate float64) string {
	feeAmount, _ := decimal.NewFromString(fee)
	rateAmount := decimal.NewFromFloat(rate)

	return feeAmount.Mul(rateAmount).Truncate(2).String()
}

// Converts and returns response payload to rate data
func toRateData(fee string, payload *exchangeRateResponsePayload) (*rateData, error) {

	sanitizedFee, err := sanitizeFee(fee)
	if err != nil {
		return nil, err
	}

	normalizedFee := normalizeFeeToBRL(sanitizedFee)

	return &rateData{
		BRL:  normalizedFee,
		USD:  toCurrencyValue(normalizedFee, payload.Rates.USD),
		EUR:  toCurrencyValue(normalizedFee, payload.Rates.EUR),
		Date: payload.Date,
	}, nil
}
