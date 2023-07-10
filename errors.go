package notify

import "fmt"

// ErrSendNotification signals that the notifier failed to send a notification.
type ErrSendNotification struct {
	Service string
	Message string
}

// Error returns the error message. It implements the error interface.
func (e *ErrSendNotification) Error() string {
	return fmt.Sprintf("failed to send notification via %s: %s", e.Service, e.Message)
}

func newErrSendNotification(service string, cause error) error {
	if cause == nil {
		return nil
	}

	return &ErrSendNotification{Service: service, Message: cause.Error()}
}

// ErrNoReceivers signals that no receivers were specified for a service.
type ErrNoReceivers struct {
	Service string
}

// Error returns the error message. It implements the error interface.
func (e *ErrNoReceivers) Error() string {
	return fmt.Sprintf("no receivers specified for service %s", e.Service)
}

// NewErrNoReceivers creates a new ErrNoReceivers error. ErrNoReceivers signals that no receivers were specified for a
// service. The service name should be passed as an argument to provide more context.
func NewErrNoReceivers(service string) error {
	return &ErrNoReceivers{Service: service}
}
