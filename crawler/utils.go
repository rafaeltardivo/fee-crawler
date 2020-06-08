package crawler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
}

// Sanitizes and returns fee amount and description
func sanitizeFeeString(rawFee string) (string, string, error) {
	description := ""

	// Tries to split text based on line breaks
	splittedValues := strings.Split(rawFee, "\n")
	lineCount := len(splittedValues)
	// Checks if values were sucessfully splitted
	hasValues := lineCount > 1
	// Checks if has description based on line count
	hasDescription := lineCount > 2

	if !hasValues {
		return "", "", crawlError("could not split amount and description")
	}

	// Ignores first line break
	feeValues := splittedValues[1:]

	// Sanitizes fee amount
	amount, err := sanitizeAmount(feeValues[0])
	if err != nil {
		logger.Error(err)
		return "", "", err
	}

	// Normalizes fee amount
	normalizedAmount := normalizeAmountToBRL(amount)

	// Sanitizes descritpion (if exists)
	if hasDescription {
		description = sanitizeDescription(feeValues[1])
	}

	return normalizedAmount, description, nil
}

// Sanitizes and returns fee amount
func sanitizeAmount(rawAmount string) (string, error) {
	re := regexp.MustCompile(`[0-9,]+[0-9]{2}`)
	fee := re.FindStringSubmatch(rawAmount)

	if len(fee) == 0 {
		return "", crawlError("fee amount not found")
	}

	amount := fee[0]
	logger.Info(fmt.Sprintf("sanitized amount: %s", amount))
	return amount, nil
}

// Sanitizes and returns fee description
func sanitizeDescription(rawDescription string) string {
	description := strings.TrimSpace(rawDescription)

	logger.Info(fmt.Sprintf("sanitized description: %s", description))
	return description
}

// Normalizes and returns the fee value using BRL floating point criteria
func normalizeAmountToBRL(fee string) string {
	return strings.Replace(fee, ",", ".", 1)
}
