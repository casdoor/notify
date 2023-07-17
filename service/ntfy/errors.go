package ntfy

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/nikoksr/notify/v2"
)

func asNotifyError(resp *http.Response) error {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Interpret the response body as an error message
	err = errors.New(strings.TrimSpace(string(b)))

	// Check the status code and return the appropriate error
	switch resp.StatusCode {
	case http.StatusOK:
		// Success
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		// Unauthorized
		return &notify.UnauthorizedError{Cause: err}
	case http.StatusTooManyRequests:
		// Rate limit
		return &notify.RateLimitError{Cause: err}
	default:
	}

	// Return a generic bad request error if none of the above matched
	return err
}
