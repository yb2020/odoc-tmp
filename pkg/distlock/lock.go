package distlock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
)

var (
	// ErrLockNotObtained is returned when a lock cannot be obtained
	ErrLockNotObtained = errors.New("lock not obtained")
	// ErrLockAlreadyReleased is returned when trying to release an already released lock
	ErrLockAlreadyReleased = errors.New("lock already released")
)

// LockOptions configures the lock behavior
type LockOptions struct {
	// Expiry is the duration for which the lock is valid
	Expiry time.Duration
	// Tries is the number of attempts to acquire the lock
	Tries int
	// RetryDelay is the delay between retry attempts
	RetryDelay time.Duration
}

// DefaultLockOptions returns the default lock options
func DefaultLockOptions() *LockOptions {
	return &LockOptions{
		Expiry:     30 * time.Second,
		Tries:      3,
		RetryDelay: 100 * time.Millisecond,
	}
}

// WithExpiry sets the expiry duration
func (o *LockOptions) WithExpiry(expiry time.Duration) *LockOptions {
	o.Expiry = expiry
	return o
}

// WithTries sets the number of tries
func (o *LockOptions) WithTries(tries int) *LockOptions {
	o.Tries = tries
	return o
}

// WithRetryDelay sets the retry delay
func (o *LockOptions) WithRetryDelay(delay time.Duration) *LockOptions {
	o.RetryDelay = delay
	return o
}

// Lock represents a distributed lock
type Lock interface {
	// Acquire attempts to acquire the lock
	Acquire(ctx context.Context) error
	// Release releases the lock
	Release(ctx context.Context) error
	// Refresh extends the lock's expiry
	Refresh(ctx context.Context, expiry time.Duration) error
	// IsAcquired returns true if the lock is currently acquired
	IsAcquired() bool
}

// Locker is the interface for creating locks
type Locker interface {
	// NewLock creates a new lock with the given key
	NewLock(key string, opts *LockOptions) Lock
	// RunWithLock runs a function with a lock
	RunWithLock(ctx context.Context, key string, opts *LockOptions, fn func(ctx context.Context) error) error
}

// RunWithLock is a helper function to run a function with a lock
func RunWithLock(ctx context.Context, locker Locker, key string, opts *LockOptions, fn func(ctx context.Context) error) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(
		ctx,
		opentracing.GlobalTracer(),
		"distlock.RunWithLock",
	)
	defer span.Finish()

	lock := locker.NewLock(key, opts)
	if err := lock.Acquire(ctx); err != nil {
		return err
	}

	defer func() {
		if err := lock.Release(ctx); err != nil && !errors.Is(err, ErrLockAlreadyReleased) {
			fmt.Printf("Failed to release lock: key=%s, error=%v\n", key, err)
		}
	}()

	return fn(ctx)
}
