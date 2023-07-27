package twilio

import (
	"errors"
	"net/http"

	twilioclient "github.com/twilio/twilio-go/client"

	"github.com/nikoksr/notify/v2"
)

func asNotifyError(err error) error {
	if err == nil {
		return nil
	}

	// If the error is not an API error, return it as is and wrap it in a bad request error
	var apiErr *twilioclient.TwilioRestError
	if !errors.As(err, &apiErr) {
		return &notify.BadRequestError{Cause: err}
	}

	switch apiErr.Code {
	case http.StatusUnauthorized, http.StatusForbidden:
		// Unauthorized
		return &notify.UnauthorizedError{Cause: err}
	case http.StatusTooManyRequests:
		// Rate limit
		return &notify.RateLimitError{Cause: err}
	}

	// If none of the above matched, return a generic bad request error
	return &notify.BadRequestError{Cause: err}
}
