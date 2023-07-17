package discord

import (
	"errors"

	"github.com/bwmarrin/discordgo"

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

	// If none of the above matched, return a generic bad request error
	return &notify.BadRequestError{Cause: err}
}
