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

// Error wrapper for external service errors
func externalServiceError(msg string) error {
	return fmt.Errorf("External service error: %s", msg)
}
