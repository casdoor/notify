package mailgun

import (
	"errors"

	"github.com/mailgun/mailgun-go/v4"

	"github.com/nikoksr/notify/v2"
	"github.com/nikoksr/notify/v2/internal/httperror"
)

func asNotifyError(err error) error {
	if err == nil {
		return nil
	}

	// Check for common errors first
	if errors.Is(err, mailgun.ErrInvalidMessage) || errors.Is(err, mailgun.ErrEmptyParam) {
		return &notify.BadRequestError{Cause: err}
	}

	// If the error is not an API error, return it as is and wrap it in a bad request error
	var apiErr *mailgun.UnexpectedResponseError
	if !errors.As(err, &apiErr) {
		return &notify.BadRequestError{Cause: err}
	}

	// If the error is an API error, use the data field as the error message
	err = errors.New(string(apiErr.Data))

	// Use the http status code to determine the appropriate Notify error
	return httperror.HandleHTTPError(err, apiErr.Actual)
}
