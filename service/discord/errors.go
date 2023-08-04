package discord

import (
	"errors"

	"github.com/bwmarrin/discordgo"

	"github.com/nikoksr/notify/v2/internal/httperror"

	"github.com/nikoksr/notify/v2"
)

func asNotifyError(err error) error {
	if err == nil {
		return nil
	}

	// Unauthorized
	if errors.Is(err, discordgo.ErrUnauthorized) {
		return &notify.UnauthorizedError{Cause: err}
	}

	// Rate limit
	var rateLimitErr *discordgo.RateLimitError
	if errors.As(err, &rateLimitErr) {
		return &notify.RateLimitError{Cause: rateLimitErr}
	}

	// If the error is not an API error, return it as is and wrap it in a bad request error
	var apiErr *discordgo.RESTError
	if !errors.As(err, &apiErr) {
		return &notify.BadRequestError{Cause: err}
	}

	// Use the http status code to determine the appropriate Notify error
	return httperror.HandleHTTPError(err, apiErr.Response.StatusCode)
}
