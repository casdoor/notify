package slack

import (
	"errors"
	"net/http"

	"github.com/slack-go/slack"

	"github.com/nikoksr/notify/v2"
)

func asNotifyError(err error) error {
	if err == nil {
		return nil
	}

	// Rate limit
	var rateLimitErr *slack.RateLimitedError
	if errors.As(err, &rateLimitErr) {
		return &notify.RateLimitError{Cause: err}
	}

	var statusCodeErr *slack.StatusCodeError
	if errors.As(err, &statusCodeErr) {
		switch statusCodeErr.Code {
		case http.StatusUnauthorized, http.StatusForbidden:
			// Unauthorized
			return &notify.UnauthorizedError{Cause: err}
		default:
		}
	}

	// If none of the above matched, return a generic bad request error
	return &notify.BadRequestError{Cause: err}
}
