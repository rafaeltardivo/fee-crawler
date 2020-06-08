package api

import "fmt"

// Error wrapper for api errors
func apiError(msg string) error {
	return fmt.Errorf("API error: %s", msg)
}
