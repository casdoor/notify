package twilio

import (
	"errors"

	twilioclient "github.com/twilio/twilio-go/client"

	"github.com/nikoksr/notify/v2/internal/httperror"

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

	// Use the http status code to determine the appropriate Notify error
	return httperror.HandleHTTPError(err, apiErr.Code)
}
