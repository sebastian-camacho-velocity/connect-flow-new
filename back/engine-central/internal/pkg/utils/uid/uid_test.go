package uid

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/xid"
)

func TestGenIdCollisions(t *testing.T) {
	// test collisions
	count := 100
	ids := make(map[string]struct{}, count)
	collisions := 0

	for i := 0; i < count; i++ {
		id := New()
		fmt.Println(id)
		if i%3 == 0 {
			time.Sleep(500 * time.Millisecond)
		}
		if _, ok := ids[id]; ok {
			collisions++
		}
		ids[id] = struct{}{}
	}

	if collisions > 0 {
		t.Errorf("collisions: %d", collisions)
	}

	// 9m4e2mr0ui3e8a215n4g
	// 1nma0215sd2wyrxq
}

func TestGenIdCollisionsParallel(t *testing.T) {
	// test collisions
	count := 10_000_000
	parts := 100

	wg := sync.WaitGroup{}
	wg.Add(parts)
	mu := sync.Mutex{}
	var collisions int32

	totalIds := map[string]struct{}{}

	for i := 0; i < parts; i++ {
		go func() {
			ids := make(map[string]struct{}, count/parts)

			for j := 0; j < count/parts; j++ {
				id := New()
				if _, ok := ids[id]; ok {
					atomic.AddInt32(&collisions, 1)
				}
				ids[id] = struct{}{}
			}

			// compare each part
			mu.Lock()
			for id := range ids {
				if _, ok := totalIds[id]; ok {
					atomic.AddInt32(&collisions, 1)
				}
				totalIds[id] = struct{}{}
			}
			mu.Unlock()

			wg.Done()
		}()
	}
	wg.Wait()

	if collisions > 0 {
		t.Errorf("collisions: %d", collisions)
	}
}

func BenchmarkGenId(b *testing.B) {
	id := ""
	for i := 0; i < b.N; i++ {
		New()
	}
	_ = id
}

func BenchmarkGenParallelId(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			New()
		}
	})
}

func BenchmarkXid(b *testing.B) {
	var id string
	for i := 0; i < b.N; i++ {
		id = xid.New().String()
	}
	_ = id
}

func BenchmarkXidParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			xid.New()
		}
	})
}

func TestXidCollitions(t *testing.T) {
	count := 10_000_00
	ids := make(map[xid.ID]struct{}, count)
	collisions := 0

	for i := 0; i < count; i++ {
		id := xid.New()
		if _, ok := ids[id]; ok {
			collisions++
		}
		ids[id] = struct{}{}
	}

	if collisions > 0 {
		t.Errorf("collisions: %d", collisions)
	}
}

func BenchmarkXxx(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandString(CharsetAlphaNum, 16)
	}
}

func BenchmarkXxx2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandString(CharsetAlphaNum, 20)
	}
}
