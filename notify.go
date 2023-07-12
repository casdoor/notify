// Package notify provides an abstraction for sending notifications
// through multiple services. It allows easy addition of new services
// and uniform handling of notification sending.
package notify

import (
	"context"
)

// Service describes a notification service that can send messages.
// Each service implementation should provide its own way of sending messages.
type Service interface {
	// Name should return a unique identifier for the service.
	Name() string
	// Send sends a message with a subject through this service.
	// Additional options can be provided to customize the sending process.
	// Returns an error if the sending process failed.
	Send(ctx context.Context, subject string, message string, opts ...SendOption) error
}

// Notify handles notifications sending through multiple services.
type Notify struct {
	services []Service // services contains all services through which notifications can be send.
}

// Option is a function that configures a Notify instance.
type Option = func(*Notify)

// WithServices adds the given services to a Notify instance's services list.
// The services will be used in the order they are provided. Nil services are ignored.
func WithServices(services ...Service) Option {
	return func(n *Notify) {
		n.UseServices(services...)
	}
}

// New creates a Notify instance with the specified options.
// If no options are provided, a Notify instance is created with no services and default options.
func New(opts ...Option) *Notify {
	n := &Notify{}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

// Create the package level Notify instance.
var defaultNotify = New()

// Default returns the standard Notify instance used by the package-level Send function.
func Default() *Notify {
	return defaultNotify
}

// UseServices appends the given service(s) to the Notify instance's services list. Nil services are ignored.
func (n *Notify) UseServices(services ...Service) {
	for _, svc := range services {
		if svc != nil {
			n.services = append(n.services, svc)
		}
	}
}

// UseServices appends the given service(s) to the defaultNotify instance's services list. Nil services are ignored.
func UseServices(services ...Service) {
	defaultNotify.UseServices(services...)
}
