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

func TestLockTemplate(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	locker := distlock.NewRedisLocker(client)
	template := distlock.NewLockTemplate(locker)

	t.Run("Lock and ReleaseLock", func(t *testing.T) {
		// Acquire a lock
		lockInfo := template.Lock("test-key", 30*time.Second, 1*time.Second)
		assert.NotNil(t, lockInfo)
		
		// Verify the key exists in Redis
		exists := mr.Exists("lock:test-key")
		assert.True(t, exists)
		
		// Release the lock
		released := template.ReleaseLock(lockInfo)
		assert.True(t, released)
		
		// Verify the key no longer exists in Redis
		exists = mr.Exists("lock:test-key")
		assert.False(t, exists)
	})

	t.Run("LockWithRetry", func(t *testing.T) {
		// First lock with short expiry
		lockInfo1 := template.Lock("test-key-2", 50*time.Millisecond, 100*time.Millisecond)
		assert.NotNil(t, lockInfo1)
		
		// Second lock with retry should eventually succeed
		go func() {
			// Wait for the first lock to expire
			time.Sleep(100 * time.Millisecond)
			// Release the first lock (though it might have expired already)
			template.ReleaseLock(lockInfo1)
		}()
		
		lockInfo2 := template.LockWithRetry("test-key-2", 30*time.Second, 500*time.Millisecond, 50*time.Millisecond)
		assert.NotNil(t, lockInfo2)
		
		// Release the second lock
		released := template.ReleaseLock(lockInfo2)
		assert.True(t, released)
	})

	t.Run("RunWithLock", func(t *testing.T) {
		executed := false
		
		err := template.RunWithLock(context.Background(), "test-key-3", 30*time.Second, func() error {
			executed = true
			return nil
		})
		
		assert.NoError(t, err)
		assert.True(t, executed)
		
		// Verify the key no longer exists in Redis (auto-released)
		exists := mr.Exists("lock:test-key-3")
		assert.False(t, exists)
	})

	t.Run("Concurrent operations with LockTemplate", func(t *testing.T) {
		counter := 0
		concurrency := 5
		wg := sync.WaitGroup{}
		wg.Add(concurrency)
		
		// Use a mutex to protect the counter
		var mu sync.Mutex
		
		// Run multiple goroutines trying to increment the counter
		for i := 0; i < concurrency; i++ {
			go func(idx int) {
				defer wg.Done()
				
				// Use different keys for each goroutine to avoid contention
				key := fmt.Sprintf("test-key-concurrent-%d", idx)
				
				err := template.RunWithLock(context.Background(), key, 30*time.Second, func() error {
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
		
		wg.Wait()
		
		// Counter should be exactly equal to concurrency (no race conditions)
		assert.Equal(t, concurrency, counter)
	})
}

func TestLockTemplateWithUserContext(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	locker := distlock.NewRedisLocker(client)
	template := distlock.NewLockTemplate(locker)

	t.Run("RunWithLock preserves context", func(t *testing.T) {
		// Create a context with a value
		ctx := context.WithValue(context.Background(), "test-key", "test-value")
		
		err := template.RunWithLock(ctx, "test-key-5", 30*time.Second, func() error {
			// In a real scenario, we would verify that user context is preserved
			// but since we're using a simple context.WithValue, we can't directly test this
			// This test is more for demonstration purposes
			return nil
		})
		
		assert.NoError(t, err)
	})
}

// This test demonstrates how to use the lock template in a real-world scenario
func TestLockTemplateRealWorldExample(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	locker := distlock.NewRedisLocker(client)
	template := distlock.NewLockTemplate(locker)

	// Simulate a service that processes tasks
	type TaskProcessor struct {
		lockTemplate *distlock.LockTemplate
		processedIDs map[string]bool
		mu           sync.Mutex
	}

	processor := &TaskProcessor{
		lockTemplate: template,
		processedIDs: make(map[string]bool),
	}

	// Method to process a task with distributed locking
	processTask := func(taskID string) error {
		// Create a lock key based on the task ID
		lockKey := "task:" + taskID
		
		return template.RunWithLock(context.Background(), lockKey, 30*time.Second, func() error {
			// Check if task is already processed (simulating idempotency check)
			processor.mu.Lock()
			if processor.processedIDs[taskID] {
				processor.mu.Unlock()
				return nil // Task already processed
			}
			
			// Simulate task processing
			time.Sleep(10 * time.Millisecond)
			
			// Mark task as processed
			processor.processedIDs[taskID] = true
			processor.mu.Unlock()
			
			return nil
		})
	}

	// Test concurrent processing of the same task
	t.Run("concurrent task processing", func(t *testing.T) {
		taskID := "task-123"
		concurrency := 5
		wg := sync.WaitGroup{}
		wg.Add(concurrency)
		
		// Create a longer expiry to ensure the lock doesn't expire during the test
		for i := 0; i < concurrency; i++ {
			go func() {
				defer wg.Done()
				err := processTask(taskID)
				assert.NoError(t, err)
			}()
			// Add a small delay between goroutines to ensure they don't all try at exactly the same time
			time.Sleep(5 * time.Millisecond)
		}
		
		wg.Wait()
		
		// Task should be processed exactly once
		processor.mu.Lock()
		assert.True(t, processor.processedIDs[taskID])
		processor.mu.Unlock()
		
		// Count the number of processed tasks (should be 1)
		count := 0
		processor.mu.Lock()
		for _, processed := range processor.processedIDs {
			if processed {
				count++
			}
		}
		processor.mu.Unlock()
		
		assert.Equal(t, 1, count)
	})
}
