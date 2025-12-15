package test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yb2020/odoc/pkg/distlock"
)

// setupRedis creates a new miniredis instance for testing
func setupRedis(t *testing.T) (*miniredis.Miniredis, redis.UniversalClient) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return mr, client
}

func TestRedisLock_Acquire(t *testing.T) {
	mr, client := setupRedis(t)
	defer mr.Close()

	locker := distlock.NewRedisLocker(client)
	
	t.Run("successful acquisition", func(t *testing.T) {
		lock := locker.NewLock("test-key", nil)
		err := lock.Acquire(context.Background())
		assert.NoError(t, err)
		assert.True(t, lock.IsAcquired())
		
		// Verify the key exists in Redis
		exists := mr.Exists("lock:test-key")
		assert.True(t, exists)
	})
	
	t.Run("failed acquisition due to existing lock", func(t *testing.T) {
		// First lock
		lock1 := locker.NewLock("test-key-2", nil)
		err := lock1.Acquire(context.Background())
		assert.NoError(t, err)
		
		// Second lock with same key should fail
		lock2 := locker.NewLock("test-key-2", nil)
		err = lock2.Acquire(context.Background())
		assert.Equal(t, distlock.ErrLockNotObtained, err)
		assert.False(t, lock2.IsAcquired())
	})
	
	t.Run("acquisition with custom options", func(t *testing.T) {
		opts := distlock.DefaultLockOptions().
			WithExpiry(100 * time.Millisecond).
			WithTries(5).
			WithRetryDelay(10 * time.Millisecond)
			
		lock := locker.NewLock("test-key-3", opts)
		err := lock.Acquire(context.Background())
		assert.NoError(t, err)
		assert.True(t, lock.IsAcquired())
		
		// Verify the key exists in Redis
		exists := mr.Exists("lock:test-key-3")
		assert.True(t, exists)
		
		// Verify TTL is close to what we set
		ttl := mr.TTL("lock:test-key-3")
		assert.Less(t, ttl, 110*time.Millisecond)
	})
	
	t.Run("acquisition with retries", func(t *testing.T) {
		// Create a lock that will expire quickly
		opts := distlock.DefaultLockOptions().
			WithExpiry(20 * time.Millisecond)
			
		lock1 := locker.NewLock("test-key-4", opts)
		err := lock1.Acquire(context.Background())
		assert.NoError(t, err)
		
		// Wait for the first lock to expire
		time.Sleep(200 * time.Millisecond)
		
		// Explicitly release the lock to ensure it's gone
		_ = lock1.Release(context.Background())
		
		// Second lock with longer retry should eventually succeed
		opts2 := distlock.DefaultLockOptions().
			WithExpiry(1 * time.Second).
			WithTries(10).
			WithRetryDelay(20 * time.Millisecond)
			
		lock2 := locker.NewLock("test-key-4", opts2)
		
		err = lock2.Acquire(context.Background())
		assert.NoError(t, err)
		assert.True(t, lock2.IsAcquired())
	})
}

func TestRedisLock_Release(t *testing.T) {
	mr, client := setupRedis(t)
	defer mr.Close()

	locker := distlock.NewRedisLocker(client)
	
	t.Run("successful release", func(t *testing.T) {
		lock := locker.NewLock("test-key", nil)
		err := lock.Acquire(context.Background())
		assert.NoError(t, err)
		
		err = lock.Release(context.Background())
		assert.NoError(t, err)
		assert.False(t, lock.IsAcquired())
		
		// Verify the key no longer exists in Redis
		exists := mr.Exists("lock:test-key")
		assert.False(t, exists)
	})
	
	t.Run("release without acquisition", func(t *testing.T) {
		lock := locker.NewLock("test-key-2", nil)
		err := lock.Release(context.Background())
		assert.Equal(t, distlock.ErrLockAlreadyReleased, err)
	})
	
	t.Run("double release", func(t *testing.T) {
		lock := locker.NewLock("test-key-3", nil)
		err := lock.Acquire(context.Background())
		assert.NoError(t, err)
		
		err = lock.Release(context.Background())
		assert.NoError(t, err)
		
		err = lock.Release(context.Background())
		assert.Equal(t, distlock.ErrLockAlreadyReleased, err)
	})
}

func TestRedisLock_Refresh(t *testing.T) {
	mr, client := setupRedis(t)
	defer mr.Close()

	locker := distlock.NewRedisLocker(client)
	
	t.Run("successful refresh", func(t *testing.T) {
		// Create a lock with short expiry
		opts := distlock.DefaultLockOptions().
			WithExpiry(100 * time.Millisecond)
			
		lock := locker.NewLock("test-key", opts)
		err := lock.Acquire(context.Background())
		assert.NoError(t, err)
		
		// Get initial TTL
		initialTTL := mr.TTL("lock:test-key")
		
		// Refresh with longer expiry
		err = lock.Refresh(context.Background(), 500*time.Millisecond)
		assert.NoError(t, err)
		
		// Verify TTL has been extended
		newTTL := mr.TTL("lock:test-key")
		assert.Greater(t, newTTL, initialTTL)
		assert.Less(t, newTTL, 550*time.Millisecond)
	})
	
	t.Run("refresh without acquisition", func(t *testing.T) {
		lock := locker.NewLock("test-key-2", nil)
		err := lock.Refresh(context.Background(), 1*time.Second)
		assert.Equal(t, distlock.ErrLockAlreadyReleased, err)
	})
	
	t.Run("refresh after release", func(t *testing.T) {
		lock := locker.NewLock("test-key-3", nil)
		err := lock.Acquire(context.Background())
		assert.NoError(t, err)
		
		err = lock.Release(context.Background())
		assert.NoError(t, err)
		
		err = lock.Refresh(context.Background(), 1*time.Second)
		assert.Equal(t, distlock.ErrLockAlreadyReleased, err)
	})
}

func TestRedisLocker_RunWithLock(t *testing.T) {
	mr, client := setupRedis(t)
	defer mr.Close()

	locker := distlock.NewRedisLocker(client)
	
	t.Run("successful execution", func(t *testing.T) {
		executed := false
		
		err := locker.RunWithLock(context.Background(), "test-key", nil, func(ctx context.Context) error {
			executed = true
			return nil
		})
		
		assert.NoError(t, err)
		assert.True(t, executed)
		
		// Verify the key no longer exists in Redis (auto-released)
		exists := mr.Exists("lock:test-key")
		assert.False(t, exists)
	})
	
	t.Run("execution with error", func(t *testing.T) {
		expectedErr := assert.AnError
		
		err := locker.RunWithLock(context.Background(), "test-key-2", nil, func(ctx context.Context) error {
			return expectedErr
		})
		
		assert.Equal(t, expectedErr, err)
		
		// Verify the key no longer exists in Redis (auto-released even on error)
		exists := mr.Exists("lock:test-key-2")
		assert.False(t, exists)
	})
	
	t.Run("concurrent execution", func(t *testing.T) {
		counter := 0
		concurrency := 5
		
		// Use a wait group to ensure all goroutines complete
		var wg sync.WaitGroup
		wg.Add(concurrency)
		
		// Use a mutex to protect the counter
		var mu sync.Mutex
		
		// Run multiple goroutines trying to increment the counter
		for i := 0; i < concurrency; i++ {
			go func(idx int) {
				defer wg.Done()
				
				// Use different keys for each goroutine to avoid contention
				key := fmt.Sprintf("test-key-concurrent-%d", idx)
				
				err := locker.RunWithLock(context.Background(), key, nil, func(ctx context.Context) error {
					mu.Lock()
					current := counter
					// Simulate some work
					time.Sleep(10 * time.Millisecond)
					counter = current + 1
					mu.Unlock()
					return nil
				})
				
				assert.NoError(t, err)
			}(i)
		}
		
		// Wait for all goroutines to complete
		wg.Wait()
		
		// Counter should be exactly equal to concurrency (no race conditions)
		assert.Equal(t, concurrency, counter)
	})
}
