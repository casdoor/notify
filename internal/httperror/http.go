package httperror

import (
	"net/http"

	"github.com/nikoksr/notify/v2"
)

func HandleHTTPError(err error, statusCode int) error {
	switch {
	case statusCode >= 400 && statusCode < 500:
		switch statusCode {
		case http.StatusUnauthorized, http.StatusForbidden:
			// Unauthorized
			return &notify.UnauthorizedError{Cause: err}
		case http.StatusTooManyRequests:
			// Rate limit
			return &notify.RateLimitError{Cause: err}
		default:
			// Other client errors (BadRequestError)
			return &notify.BadRequestError{Cause: err}
		}
	case statusCode >= 500:
		// ExternalServerError for 5xx status codes
		return &notify.ExternalServerError{Cause: err}
	default:
		// Other status codes
		return err
	}
}
