package notify

import (
	"errors"
	"fmt"
)

// ErrNoReceivers indicates that there are no receivers specified for a service.
var ErrNoReceivers = errors.New("no receivers specified")

// ErrSendNotification encapsulates any errors that occur when sending a notification.
type ErrSendNotification struct {
	// ReceiverID is the ID of the receiver that failed to receive the notification.
	ReceiverID any
	// Cause is the underlying error that caused the notification to fail.
	Cause error
}

// Error provides the string representation of the ErrSendNotification error.
func (e *ErrSendNotification) Error() string {
	return fmt.Sprintf("Failed to send notification to receiver: %v, cause: %v", e.ReceiverID, e.Cause)
}

// Unwrap retrieves the underlying error for the ErrSendNotification error.
func (e *ErrSendNotification) Unwrap() error {
	return e.Cause
}

// NewErrSendNotification is a factory function that creates and returns a new ErrSendNotification error.
func NewErrSendNotification(receiver any, cause error) *ErrSendNotification {
	return &ErrSendNotification{
		ReceiverID: receiver,
		Cause:      cause,
	}
}

// ErrServiceFailure represents an error that occurs when a service fails.
type ErrServiceFailure struct {
	// Service is the name of the service that failed.
	Service string
	// Err is the underlying error that caused the service to fail.
	Err error
}

// Unwrap retrieves the underlying error for the ErrServiceFailure error.
func (e *ErrServiceFailure) Unwrap() error {
	return e.Err
}

// Error provides the string representation of the ErrServiceFailure error.
func (e *ErrServiceFailure) Error() string {
	return fmt.Sprintf("%s: %s", e.Service, e.Err)
}

// newErrServiceFailure is a factory function that creates and returns a new ErrServiceFailure error.
func newErrServiceFailure(service string, err error) *ErrServiceFailure {
	return &ErrServiceFailure{
		Service: service,
		Err:     err,
	}
}
