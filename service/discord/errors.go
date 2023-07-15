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

	if errors.Is(err, discordgo.ErrUnauthorized) {
		return &notify.UnauthorizedError{Cause: err}
	}

	var rateLimitErr *discordgo.RateLimitError
	if errors.As(err, &rateLimitErr) {
		return &notify.RateLimitError{Cause: rateLimitErr}
	}

	return err
}
