package test

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/distlock"
)

// This test demonstrates how to use the distributed lock in a real-world scenario
// that's similar to the full-text translation service
func TestDistLockWithTranslationService(t *testing.T) {
	// Set up a mini Redis server for testing
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Create the lock components
	locker := distlock.NewRedisLocker(client)
	template := distlock.NewLockTemplate(locker)

	// Simulate a translation service
	type TranslationService struct {
		lockTemplate *distlock.LockTemplate
		translations map[int64]string      // pdfId -> translated content
		inProgress   map[int64]bool        // pdfId -> is translation in progress
		mu           sync.Mutex
	}

	service := &TranslationService{
		lockTemplate: template,
		translations: make(map[int64]string),
		inProgress:   make(map[int64]bool),
	}

	// Method to translate a PDF with distributed locking
	translatePDF := func(ctx context.Context, userID, pdfID int64) (string, error) {
		// Create a lock key based on the PDF ID (similar to what would be done in the real service)
		lockKey := "translate:pdf:" + strconv.FormatInt(pdfID, 10)
		
		// Create a user context for the background operations
		uc := userContext.NewUserContext().SetUserID(userID)
		ctx = uc.ToContext(ctx)
		
		// First check if translation already exists (without lock)
		service.mu.Lock()
		if content, exists := service.translations[pdfID]; exists {
			service.mu.Unlock()
			return content, nil
		}
		service.mu.Unlock()
		
		// Use the lock template to ensure only one translation process runs at a time
		var translationContent string
		
		err := template.RunWithLock(ctx, lockKey, 30*time.Second, func() error {
			// Double-check after acquiring the lock
			service.mu.Lock()
			if content, exists := service.translations[pdfID]; exists {
				service.mu.Unlock()
				translationContent = content
				return nil
			}
			
			// Check if translation is already in progress
			if service.inProgress[pdfID] {
				service.mu.Unlock()
				translationContent = "Translation in progress"
				return nil
			}
			
			// Mark as in progress
			service.inProgress[pdfID] = true
			service.mu.Unlock()
			
			// Simulate async translation process
			go func() {
				// Create a background context that won't be canceled
				bgCtx := context.Background()
				
				// If we have user context, preserve it
				uc := userContext.GetUserContext(ctx)
				if uc != nil {
					bgCtx = uc.ToContext(bgCtx)
				}
				
				// Simulate translation work
				time.Sleep(100 * time.Millisecond)
				
				// Store the result
				service.mu.Lock()
				service.translations[pdfID] = "Translated content for PDF " + strconv.FormatInt(pdfID, 10)
				service.inProgress[pdfID] = false
				service.mu.Unlock()
			}()
			
			translationContent = "Translation started"
			return nil
		})
		
		return translationContent, err
	}

	// Test concurrent translation requests for the same PDF
	t.Run("concurrent translation requests", func(t *testing.T) {
		pdfID := int64(123)
		userID := int64(456)
		concurrency := 5
		results := make([]string, concurrency)
		
		// Create a wait group to synchronize the test
		wg := sync.WaitGroup{}
		wg.Add(concurrency)
		
		// Launch multiple concurrent translation requests
		for i := 0; i < concurrency; i++ {
			go func(index int) {
				defer wg.Done()
				result, err := translatePDF(context.Background(), userID, pdfID)
				assert.NoError(t, err)
				results[index] = result
			}(i)
			// Add a small delay between goroutines to ensure they don't all try at exactly the same time
			time.Sleep(5 * time.Millisecond)
		}
		
		wg.Wait()
		
		// Verify that all results are valid responses
		for i := 0; i < concurrency; i++ {
			assert.Contains(t, []string{"Translation started", "Translation in progress"}, results[i],
				"Each result should be either 'Translation started' or 'Translation in progress'")
		}
		
		// Wait for the background translation to complete
		time.Sleep(200 * time.Millisecond)
		
		// Verify the translation was completed
		service.mu.Lock()
		assert.Contains(t, service.translations, pdfID)
		assert.False(t, service.inProgress[pdfID])
		service.mu.Unlock()
		
		// A new request should now get the completed translation
		result, err := translatePDF(context.Background(), userID, pdfID)
		assert.NoError(t, err)
		assert.Equal(t, "Translated content for PDF "+strconv.FormatInt(pdfID, 10), result)
	})
}
