package ntfy

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/nikoksr/notify/v2"
)

func asNotifyError(resp *http.Response) error {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = errors.New(strings.TrimSpace(string(b)))

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return &notify.UnauthorizedError{Cause: err}
	case http.StatusTooManyRequests:
		return &notify.RateLimitError{Cause: err}
	default:
	}

	return err
}
