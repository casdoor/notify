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

	var apiErr *telegram.Error
	if !errors.As(err, &apiErr) {
		return err
	}

	switch apiErr.Code {
	case http.StatusUnauthorized, http.StatusForbidden:
		return &notify.UnauthorizedError{Cause: err}
	case http.StatusTooManyRequests:
		return &notify.RateLimitError{Cause: err}
	}

	return err
}
