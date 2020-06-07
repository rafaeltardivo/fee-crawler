package rates

import (
	"fmt"
)

// Error wrapper for parse errors
func parseError(msg string) error {
	return fmt.Errorf("Parse error: %s", msg)
}

// Error wrapper for database errors
func databaseError(msg string) error {
	return fmt.Errorf("Database error: %s", msg)
}

// Error wrapper for api errors
func apiError(msg string) error {
	return fmt.Errorf("API error: %s", msg)
}
