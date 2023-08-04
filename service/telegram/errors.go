package telegram

import (
	"errors"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/nikoksr/notify/v2/internal/httperror"

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

	// Use the http status code to determine the appropriate Notify error
	return httperror.HandleHTTPError(err, apiErr.Code)
}
