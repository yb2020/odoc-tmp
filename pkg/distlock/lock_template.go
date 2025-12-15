package distlock

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	userContext "github.com/yb2020/odoc/pkg/context"
)

// LockTemplate provides a higher-level API for distributed locking
type LockTemplate struct {
	locker Locker
}

// NewLockTemplate creates a new LockTemplate
func NewLockTemplate(locker Locker) *LockTemplate {
	return &LockTemplate{
		locker: locker,
	}
}

// Lock acquires a lock and returns a LockInfo
func (t *LockTemplate) Lock(key string, expiry time.Duration, timeout time.Duration) *LockInfo {
	return t.LockWithRetry(key, expiry, timeout, 100*time.Millisecond)
}

// LockWithRetry acquires a lock with retry parameters and returns a LockInfo
func (t *LockTemplate) LockWithRetry(key string, expiry time.Duration, timeout time.Duration, retryDelay time.Duration) *LockInfo {
	// Create a background context that won't be canceled
	ctx := context.Background()

	// Calculate number of tries based on timeout and retry delay
	tries := int(timeout / retryDelay)
	if tries < 1 {
		tries = 1
	}

	opts := DefaultLockOptions().
		WithExpiry(expiry).
		WithTries(tries).
		WithRetryDelay(retryDelay)

	lock := t.locker.NewLock(key, opts)
	err := lock.Acquire(ctx)
	if err != nil {
		return nil
	}

	return &LockInfo{
		Key:    key,
		Lock:   lock,
		Expiry: expiry,
	}
}

// RunWithLock runs a function with a lock, using a background context
// to ensure the lock operation completes even if the original context is canceled
func (t *LockTemplate) RunWithLock(ctx context.Context, key string, expiry time.Duration, fn func() error) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(
		ctx,
		opentracing.GlobalTracer(),
		"LockTemplate.RunWithLock",
	)
	defer span.Finish()

	// Create a background context that won't be canceled
	bgCtx := context.Background()

	// If we have user context, preserve it
	uc := userContext.GetUserContext(ctx)
	if uc != nil {
		bgCtx = uc.ToContext(bgCtx)
	}

	opts := DefaultLockOptions().WithExpiry(expiry)

	// Adapt the function to accept a context parameter
	return t.locker.RunWithLock(bgCtx, key, opts, func(lockCtx context.Context) error {
		return fn()
	})
}

// LockInfo contains information about an acquired lock
type LockInfo struct {
	Key    string
	Lock   Lock
	Expiry time.Duration
}

// ReleaseLock releases the lock
func (t *LockTemplate) ReleaseLock(lockInfo *LockInfo) bool {
	if lockInfo == nil || lockInfo.Lock == nil {
		return false
	}

	// Create a background context that won't be canceled
	ctx := context.Background()

	err := lockInfo.Lock.Release(ctx)
	return err == nil
}
