package retry

import (
	"errors"
	"math/rand"
	"time"
)

// ExponentialBackoff contains the configuration for exponential backoff
type ExponentialBackoff struct {
	// InitialDelayInterval represents the initial amount of time of sleep
	InitialDelayInterval time.Duration
	// MaxDelayInterval represents the maximal amount of time of sleep between
	// retries
	MaxDelayInterval time.Duration
	// MaxRetryAttempts is the maximal number of retries before exiting with
	// an error. A value of zero signifies unlimited retry attemps
	MaxRetryAttempts int
	// Multiplier is used to increment the current interval by multiplying it with
	// this multiplier
	Multiplier float64
}

// Func represents a function to retry, which returns true if the attempt is
// successful, or an error if the retry should be aborted
type Func func(retry int) (done bool, err error)

// ErrMaxRetryAttempts is returned when the number of maximal retry attempts is
// reached
var ErrMaxRetryAttempts = errors.New("maximal number of retry attempts reached")

// Retry retries the provided func with exponential backoff, until
// the maximal number of retries is reached
func (b *ExponentialBackoff) Retry(fn Func) error {
	wait := b.InitialDelayInterval

	for i := 0; i < b.MaxRetryAttempts || b.MaxRetryAttempts == 0; i++ {
		if i != 0 {
			// Sleep for the determined duration
			time.Sleep(wait)

			// Exponentially increase that sleep duration
			wait = time.Duration(float64(wait) * b.Multiplier)

			// Add a jitter (randomized delay) for the next attempt, to prevent
			// potential collisions
			wait = wait + time.Duration(rand.Float64()*float64(wait))
		}

		if ok, err := fn(i); err != nil || ok {
			return err
		}
	}

	return ErrMaxRetryAttempts
}
