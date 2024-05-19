package utils

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParallelize(t *testing.T) {
	t.Run("Should run functions in parallel and wait for all to complete", func(t *testing.T) {
		var mu sync.Mutex
		results := []int{}

		fn1 := func() {
			time.Sleep(100 * time.Millisecond)
			mu.Lock()
			results = append(results, 1)
			mu.Unlock()
		}

		fn2 := func() {
			time.Sleep(50 * time.Millisecond)
			mu.Lock()
			results = append(results, 2)
			mu.Unlock()
		}

		start := time.Now()
		Parallelize(fn1, fn2)
		duration := time.Since(start)

		assert.ElementsMatch(t, []int{1, 2}, results)
		assert.Less(t, duration.Milliseconds(), int64(150))
	})

	t.Run("Should handle no functions gracefully", func(t *testing.T) {
		start := time.Now()
		Parallelize()
		duration := time.Since(start)

		assert.Less(t, duration.Milliseconds(), int64(50))
	})

	t.Run("Should run a single function", func(t *testing.T) {
		var result int

		fn := func() {
			result = 1
		}

		Parallelize(fn)

		assert.Equal(t, 1, result)
	})
}
