package telegram

import (
	"errors"
	"net/http"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/nikoksr/notify/v2"
)

func asNotifyError(err error) error {
	if err == nil {
		return nil
	}

	// If the error is not an API error, return it as is and wrap it in a bad request error
	var apiErr *telegram.Error
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
