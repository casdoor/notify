package sendgrid

import (
	"errors"
	"fmt"

	"github.com/sendgrid/rest"

	"github.com/nikoksr/notify/v2/internal/httperror"
)

func asNotifyError(resp *rest.Response, err error) error {
	if resp == nil {
		return err
	}

	// Interpret the response body as an error message
	if resp.Body != "" {
		err = fmt.Errorf("%s: %w", err, errors.New(resp.Body))
	}

	// Use the http status code to determine the appropriate Notify error
	return httperror.HandleHTTPError(err, resp.StatusCode)
}
