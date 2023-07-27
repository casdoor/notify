// Package notify provides an abstraction for sending notifications
// through multiple services. It allows easy addition of new services
// and uniform handling of notification sending.
package notify

import "context"

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

// Dispatcher handles notifications sending through multiple services.
type Dispatcher struct {
	services []Service // services contains all services through which notifications can be send.
}

// Option is a function that configures a Dispatcher instance.
type Option = func(*Dispatcher)

// WithServices adds the given services to a Dispatcher instance's services list.
// The services will be used in the order they are provided. Nil services are ignored.
func WithServices(services ...Service) Option {
	return func(d *Dispatcher) {
		d.UseServices(services...)
	}
}

// New creates a Dispatcher instance with the specified options.
// If no options are provided, a Dispatcher instance is created with no services and default options.
func New(opts ...Option) *Dispatcher {
	n := &Dispatcher{}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

// Create the package level Dispatcher instance.
var defaultDispatcher = New()

// Default returns the standard Dispatcher instance used by the package-level Send function.
func Default() *Dispatcher {
	return defaultDispatcher
}

// SetDefault sets the package-level Dispatcher instance to the given Dispatcher instance. Nil values are ignored.
func SetDefault(d *Dispatcher) {
	if d != nil {
		defaultDispatcher = d
	}
}

// UseServices appends the given service(s) to the Dispatcher instance's services list. Nil services are ignored.
func (d *Dispatcher) UseServices(services ...Service) {
	for _, svc := range services {
		if svc != nil {
			d.services = append(d.services, svc)
		}
	}
}

// UseServices appends the given service(s) to the defaultNotify instance's services list. Nil services are ignored.
func UseServices(services ...Service) {
	defaultDispatcher.UseServices(services...)
}
