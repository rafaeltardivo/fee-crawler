package crawler

import (
	"fmt"
)

// Error wrapper for crawl errors
func crawlError(msg string) error {
	return fmt.Errorf("Crawler error: %s", msg)
}
