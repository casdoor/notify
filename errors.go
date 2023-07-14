package notify

import (
	"errors"
	"fmt"
)

// ErrNoRecipients indicates that there are no recipients specified for a service.
var ErrNoRecipients = errors.New("no recipients specified")

// ErrUnauthorized indicates that the user is not authorized to perform the requested action.
type ErrUnauthorized struct {
	// Cause is the underlying error that caused the unauthorized error.
	Cause error
}

// Error provides the string representation of the ErrUnauthorized error.
func (e *ErrUnauthorized) Error() string {
	return fmt.Sprintf("unauthorized: %v", e.Cause)
}

// Unwrap retrieves the underlying error for the ErrUnauthorized error.
func (e *ErrUnauthorized) Unwrap() error {
	return e.Cause
}

// ErrRateLimitExceeded indicates that the rate limit for the service has been exceeded.
type ErrRateLimitExceeded struct {
	// Cause is the underlying error that caused the rate limit exceeded error.
	Cause error
}

// Error provides the string representation of the ErrRateLimitExceeded error.
func (e *ErrRateLimitExceeded) Error() string {
	return fmt.Sprintf("rate limit exceeded: %v", e.Cause)
}

// Unwrap retrieves the underlying error for the ErrRateLimitExceeded error.
func (e *ErrRateLimitExceeded) Unwrap() error {
	return e.Cause
}

// ErrSendNotification encapsulates any errors that occur when sending a notification.
type ErrSendNotification struct {
	// Recipient is the identifier of the recipient that the notification failed to send to. Commonly an email address,
	// phone number, user ID, or a webhook URL.
	Recipient any
	// Cause is the underlying error that caused the notification to fail.
	Cause error
}

// Error provides the string representation of the ErrSendNotification error.
func (e *ErrSendNotification) Error() string {
	return fmt.Sprintf("failed to send notification to recipient %v: %v", e.Recipient, e.Cause)
}

// Unwrap retrieves the underlying error for the ErrSendNotification error.
func (e *ErrSendNotification) Unwrap() error {
	return e.Cause
}

// ErrServiceFailure represents an error that occurs when a service fails.
type ErrServiceFailure struct {
	// Service is the name of the service that failed.
	Service string
	// Cause is the underlying error that caused the service to fail.
	Cause error
}

// Unwrap retrieves the underlying error for the ErrServiceFailure error.
func (e *ErrServiceFailure) Unwrap() error {
	return e.Cause
}

// Error provides the string representation of the ErrServiceFailure error.
func (e *ErrServiceFailure) Error() string {
	return fmt.Sprintf("%s: %s", e.Service, e.Cause)
}
