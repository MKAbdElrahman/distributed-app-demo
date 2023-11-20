package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

// HTTPSender is an implementation of the Sender interface using HTTP.
type HTTPLogger struct {
	Endpoint string
}

// Send sends the given message to the specified endpoint using an HTTP POST request.
func (hs *HTTPLogger) Log(message string) error {
	resp, err := http.Post(hs.Endpoint, "application/text", strings.NewReader(message))
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
