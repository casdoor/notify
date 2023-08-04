package slack

import (
	"errors"
	"fmt"

	"github.com/slack-go/slack"

	"github.com/nikoksr/notify/v2/internal/httperror"

	"github.com/nikoksr/notify/v2"
)

func slackErrorResponseToError(err *slack.SlackErrorResponse) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	for _, message := range err.ResponseMetadata.Warnings {
		errMsg += fmt.Sprintf(": %s", message)
	}

	return errors.New(errMsg)
}

func asNotifyError(err error) error {
	if err == nil {
		return nil
	}

	// Rate limit
	var rateLimitErr *slack.RateLimitedError
	if errors.As(err, &rateLimitErr) {
		return &notify.RateLimitError{Cause: err}
	}

	// Check for common errors first

	var httpErr *slack.StatusCodeError
	if errors.As(err, &httpErr) {
		err = errors.New(httpErr.Error())

		// Use the http status code to determine the appropriate Notify error
		return httperror.HandleHTTPError(err, httpErr.Code)
	}

	var apiErr *slack.SlackErrorResponse
	if errors.As(err, &apiErr) {
		return slackErrorResponseToError(apiErr)
	}

	// Unable to determine error type
	return &notify.BadRequestError{Cause: err}
}
