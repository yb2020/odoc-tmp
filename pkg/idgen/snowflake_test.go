package idgen

import (
	"sync"
	"testing"
	"time"
)

func TestSnowflakeIDGenerator_NextID(t *testing.T) {
	generator, err := NewSnowflakeIDGenerator()
	if err != nil {
		t.Fatalf("Failed to create snowflake generator: %v", err)
	}

	// 测试基本ID生成
	id, err := generator.NextID()
	if err != nil {
		t.Fatalf("Failed to generate ID: %v", err)
	}
	if id <= 0 {
		t.Error("Generated ID should be positive")
	}

	// 测试时间戳部分
	timestamp := (uint64(id) >> timestampLeftShift) + twepoch
	now := uint64(time.Now().UnixNano() / 1e6)
	if timestamp > now {
		t.Errorf("Generated timestamp %d is in the future (now: %d)", timestamp, now)
	}
	if now-timestamp > 1000 { // 允许1秒的误差
		t.Errorf("Generated timestamp %d is too old (now: %d)", timestamp, now)
	}

	// 测试数据中心ID部分
	dataCenterId := (uint64(id) >> dataCenterIdShift) & maxDataCenterId
	if dataCenterId > maxDataCenterId {
		t.Errorf("Data center ID %d exceeds maximum %d", dataCenterId, maxDataCenterId)
	}

	// 测试序列号部分
	sequence := uint64(id) & sequenceMask
	if sequence > sequenceMask {
		t.Errorf("Sequence %d exceeds maximum %d", sequence, sequenceMask)
	}

	// 测试ID的唯一性和单调递增性
	t.Run("Uniqueness and Monotonicity", func(t *testing.T) {
		const n = 1000
		ids := make([]int64, n)
		for i := 0; i < n; i++ {
			id, err := generator.NextID()
			if err != nil {
				t.Fatalf("Failed to generate ID at iteration %d: %v", i, err)
			}
			ids[i] = id
		}

		// 检查单调递增性
		for i := 1; i < n; i++ {
			if ids[i] <= ids[i-1] {
				t.Errorf("IDs not monotonically increasing at index %d: %d <= %d", i, ids[i], ids[i-1])
			}
		}

		// 检查唯一性
		seen := make(map[int64]bool)
		for i, id := range ids {
			if seen[id] {
				t.Errorf("Duplicate ID found at index %d: %d", i, id)
			}
			seen[id] = true
		}
	})

	// 测试并发安全性
	t.Run("Concurrency Safety", func(t *testing.T) {
		const (
			numGoroutines = 10
			idsPerRoutine = 1000
		)
		var wg sync.WaitGroup
		ids := make(chan int64, numGoroutines*idsPerRoutine)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < idsPerRoutine; j++ {
					id, err := generator.NextID()
					if err != nil {
						t.Errorf("Failed to generate ID: %v", err)
						return
					}
					ids <- id
				}
			}()
		}

		wg.Wait()
		close(ids)

		// 检查所有生成的ID的唯一性
		seen := make(map[int64]bool)
		for id := range ids {
			if seen[id] {
				t.Errorf("Duplicate ID found in concurrent test: %d", id)
			}
			seen[id] = true
		}

		if len(seen) != numGoroutines*idsPerRoutine {
			t.Errorf("Expected %d unique IDs, got %d", numGoroutines*idsPerRoutine, len(seen))
		}
	})
}

func TestSnowflakeIDGenerator_ClockBackwards(t *testing.T) {
	generator, err := NewSnowflakeIDGenerator()
	if err != nil {
		t.Fatalf("Failed to create snowflake generator: %v", err)
	}

	// 获取第一个ID
	_, err = generator.NextID()
	if err != nil {
		t.Fatalf("Failed to generate first ID: %v", err)
	}

	// 手动设置上一次时间戳为未来时间
	generator.lastTimestamp = uint64(time.Now().UnixNano()/1e6) + 1000

	// 尝试获取下一个ID，应该返回错误
	_, err = generator.NextID()
	if err == nil {
		t.Error("Expected error when clock moves backwards, got nil")
	}
	if err.Error() != "clock moved backwards, refusing to generate id" {
		t.Errorf("Expected clock backwards error, got: %v", err)
	}
}
