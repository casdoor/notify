package ntfy

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/nikoksr/notify/v2/internal/httperror"
)

func asNotifyError(resp *http.Response) error {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Interpret the response body as an error message
	err = errors.New(strings.TrimSpace(string(b)))

	// Use the http status code to determine the appropriate Notify error
	return httperror.HandleHTTPError(err, resp.StatusCode)
}
