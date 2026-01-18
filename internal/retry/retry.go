
package retry

import (
	"time"
)

// Func is a function that can be retried.
type Func func() error

// Do retries a function with exponential backoff.
func Do(fn Func, maxRetries int, initialDelay time.Duration) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(initialDelay)
		initialDelay *= 2
	}
	return err
}
