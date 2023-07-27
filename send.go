package notify

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// Send sends a notification with the given subject and message through all the services of n. It performs these
// operations concurrently and returns the first encountered error, if any.
func (n *Notify) Send(ctx context.Context, subject, message string, opts ...SendOption) error {
	var eg errgroup.Group
	for _, service := range n.services {
		service := service

		eg.Go(func() error {
			if err := service.Send(ctx, subject, message, opts...); err != nil {
				return &ServiceFailureError{Service: service.Name(), Cause: err}
			}

			return nil
		})
	}

	return eg.Wait()
}

// Send sends a notification with the given subject and message through all the services of the defaultNotify instance.
// It performs these operations concurrently and returns the first encountered error, if any.
func Send(ctx context.Context, subject, message string, opts ...SendOption) error {
	return defaultNotify.Send(ctx, subject, message, opts...)
}
