package main

import (
	"sync"
	"testing"
)

func TestConcurrentAccess(t *testing.T) {
	// 初始化 Redis 连接
	InitRedis()

	key := "test-key"
	value := "test-value"
	const goroutineCount = 10

	// 设置初始值
	err := SetValueWithRandomExp(Clients.Test, key, value)
	if err != nil {
		t.Fatalf("Error setting value: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func() {
			defer wg.Done()

			val, err := GetKeyAndUpdateExpiration(Clients.Test, key)
			if err != nil {
				t.Errorf("Error getting value: %v", err)
			} else if val != value {
				t.Errorf("Expected value %s, got %s", value, val)
			}
		}()
	}

	wg.Wait()

	accessCount := keyAccessMap[key]
	if accessCount != goroutineCount {
		t.Errorf("Expected access count %d, got %d", goroutineCount, accessCount)
	}
}
