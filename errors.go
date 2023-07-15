package notify

import (
	"errors"
	"fmt"
)

// ErrNoRecipients indicates that there are no recipients specified for a service.
var ErrNoRecipients = errors.New("no recipients specified")

// UnauthorizedError indicates that the user is not authorized to perform the requested action.
type UnauthorizedError struct {
	// Cause is the underlying error that caused the unauthorized error.
	Cause error
}

// Error provides the string representation of the UnauthorizedError error.
func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized: %v", e.Cause)
}

// Unwrap retrieves the underlying error for the UnauthorizedError error.
func (e *UnauthorizedError) Unwrap() error {
	return e.Cause
}

// RateLimitError indicates that the rate limit for the service has been exceeded.
type RateLimitError struct {
	// Cause is the underlying error that caused the rate limit exceeded error.
	Cause error
}

// Error provides the string representation of the RateLimitError error.
func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limit exceeded: %v", e.Cause)
}

// Unwrap retrieves the underlying error for the RateLimitError error.
func (e *RateLimitError) Unwrap() error {
	return e.Cause
}

// SendNotificationError encapsulates any errors that occur when sending a notification.
type SendNotificationError struct {
	// Recipient is the identifier of the recipient that the notification failed to send to. Commonly an email address,
	// phone number, user ID, or a webhook URL.
	Recipient any
	// Cause is the underlying error that caused the notification to fail.
	Cause error
}

// Error provides the string representation of the SendNotificationError error.
func (e *SendNotificationError) Error() string {
	return fmt.Sprintf("failed to send notification to recipient %v: %v", e.Recipient, e.Cause)
}

// Unwrap retrieves the underlying error for the SendNotificationError error.
func (e *SendNotificationError) Unwrap() error {
	return e.Cause
}

// ServiceFailureError represents an error that occurs when a service fails.
type ServiceFailureError struct {
	// Service is the name of the service that failed.
	Service string
	// Cause is the underlying error that caused the service to fail.
	Cause error
}

// Unwrap retrieves the underlying error for the ServiceFailureError error.
func (e *ServiceFailureError) Unwrap() error {
	return e.Cause
}

// Error provides the string representation of the ServiceFailureError error.
func (e *ServiceFailureError) Error() string {
	return fmt.Sprintf("%s: %s", e.Service, e.Cause)
}
