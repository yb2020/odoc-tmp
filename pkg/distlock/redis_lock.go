package distlock

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/internal/database"
	userContext "github.com/yb2020/odoc/pkg/context"
)

const (
	lockScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("PEXPIRE", KEYS[1], ARGV[2])
else
    return 0
end`

	unlockScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
)

// RedisLock implements a distributed lock using Redis
type RedisLock struct {
	client  database.RedisClient
	key     string
	value   string
	expiry  time.Duration
	options *LockOptions
	locked  bool
}

// RedisLocker creates Redis-based locks
type RedisLocker struct {
	client database.RedisClient
}

// NewRedisLocker creates a new Redis-based locker
func NewRedisLocker(client database.RedisClient) *RedisLocker {
	return &RedisLocker{
		client: client,
	}
}

// NewLock creates a new Redis-based lock
func (l *RedisLocker) NewLock(key string, opts *LockOptions) Lock {
	if opts == nil {
		opts = DefaultLockOptions()
	}

	// Generate a random value to identify this lock instance
	value := make([]byte, 16)
	rand.Read(value)

	return &RedisLock{
		client:  l.client,
		key:     "lock:" + key,
		value:   hex.EncodeToString(value),
		expiry:  opts.Expiry,
		options: opts,
	}
}

// RunWithLock runs a function with a Redis lock
func (l *RedisLocker) RunWithLock(ctx context.Context, key string, opts *LockOptions, fn func(ctx context.Context) error) error {
	return RunWithLock(ctx, l, key, opts, fn)
}

// Acquire attempts to acquire the lock
func (r *RedisLock) Acquire(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(
		ctx,
		opentracing.GlobalTracer(),
		"RedisLock.Acquire",
	)
	defer span.Finish()

	// Create a background context to ensure the lock operation completes
	// even if the original context is canceled
	bgCtx := context.Background()

	// If we have user context, preserve it
	uc := userContext.GetUserContext(ctx)
	if uc != nil {
		bgCtx = uc.ToContext(bgCtx)
	}

	for i := 0; i < r.options.Tries; i++ {
		ok, err := r.client.SetNX(bgCtx, r.key, r.value, r.expiry).Result()
		if err != nil {
			return err
		}
		if ok {
			r.locked = true
			return nil
		}

		// If this is not the last try, sleep before retrying
		if i < r.options.Tries-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(r.options.RetryDelay):
			}
		}
	}

	return ErrLockNotObtained
}

// Release releases the lock
func (r *RedisLock) Release(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(
		ctx,
		opentracing.GlobalTracer(),
		"RedisLock.Release",
	)
	defer span.Finish()

	if !r.locked {
		return ErrLockAlreadyReleased
	}

	// Create a background context to ensure the lock operation completes
	// even if the original context is canceled
	bgCtx := context.Background()

	// If we have user context, preserve it
	uc := userContext.GetUserContext(ctx)
	if uc != nil {
		bgCtx = uc.ToContext(bgCtx)
	}

	res, err := r.client.Eval(bgCtx, unlockScript, []string{r.key}, r.value).Int64()
	if err != nil {
		return err
	}

	r.locked = false

	if res == 0 {
		return ErrLockAlreadyReleased
	}

	return nil
}

// Refresh extends the lock's expiry
func (r *RedisLock) Refresh(ctx context.Context, expiry time.Duration) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(
		ctx,
		opentracing.GlobalTracer(),
		"RedisLock.Refresh",
	)
	defer span.Finish()

	if !r.locked {
		return ErrLockAlreadyReleased
	}

	// Create a background context to ensure the lock operation completes
	// even if the original context is canceled
	bgCtx := context.Background()

	// If we have user context, preserve it
	uc := userContext.GetUserContext(ctx)
	if uc != nil {
		bgCtx = uc.ToContext(bgCtx)
	}

	ttlMillis := int64(expiry / time.Millisecond)
	res, err := r.client.Eval(bgCtx, lockScript, []string{r.key}, r.value, ttlMillis).Int64()
	if err != nil {
		return err
	}

	if res == 0 {
		r.locked = false
		return errors.New("lock lost")
	}

	r.expiry = expiry
	return nil
}

// IsAcquired returns true if the lock is currently acquired
func (r *RedisLock) IsAcquired() bool {
	return r.locked
}
