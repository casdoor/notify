package notify

import (
	"context"
)

// Service defines the behavior for notification services.
//
// Name returns the name of the service. This is used to identify the service in error messages.
//
// Send sends a message to the service. It returns an error if the message could not be sent. You can use the SendOption
// functions to configure the send behavior.
type Service interface {
	Name() string
	Send(ctx context.Context, subject string, message string, opts ...SendOption) error
}

// Notify is the central struct for managing notification services and sending messages to them.
type Notify struct {
	services []Service
}

// Option is a function that configures a Notify instance.
type Option = func(*Notify)

// WithServices configures the Notify instance to use the specified services. The given services are used in the order
// they are provided and are appended to any existing services. nil services are ignored.
func WithServices(services ...Service) Option {
	return func(n *Notify) {
		n.UseServices(services...)
	}
}

// New creates a new Notify instance with the given options. If no options are provided, the Notify instance is created
// with no services and default options.
func New(opts ...Option) *Notify {
	n := &Notify{}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

// Create the package level Notify instance.
var defaultNotify = New()

// Default returns the standard Notify instance used by the package-level send function.
func Default() *Notify {
	return defaultNotify
}

// useServices adds the given service(s) to the Service's services list.
func (n *Notify) useServices(services ...Service) {
	for _, svc := range services {
		if svc != nil {
			n.services = append(n.services, svc)
		}
	}
}

// UseServices adds the given service(s) to the Service's services list.
func (n *Notify) UseServices(services ...Service) {
	n.useServices(services...)
}

// UseServices adds the given service(s) to the Service's services list.
func UseServices(services ...Service) {
	defaultNotify.UseServices(services...)
}
