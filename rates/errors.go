package rates

import (
	"fmt"
)

// Error wrapper for parse errors
func parseError(msg string) error {
	return fmt.Errorf("Parse error: %s", msg)
}

// Error wrapper for client errors
func clientError(msg string) error {
	return fmt.Errorf("Client error: %s", msg)
}

// Error wrapper for convert errors
func convertError(msg string) error {
	return fmt.Errorf("Convert error: %s", msg)
}
