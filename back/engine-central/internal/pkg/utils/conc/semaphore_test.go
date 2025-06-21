package conc

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestNewSemaphore(t *testing.T) {
	concurrency := 3
	sem := NewSemaphore(concurrency)
	if cap(sem.ch) != concurrency {
		t.Errorf("expected channel capacity %d, got %d", concurrency, cap(sem.ch))
	}
}

func TestSemaphoreRun(t *testing.T) {
	concurrency := 2
	sem := NewSemaphore(concurrency)
	var running int32

	task := func() {
		atomic.AddInt32(&running, 1)
		time.Sleep(100 * time.Millisecond)
		atomic.AddInt32(&running, -1)
	}

	for i := 0; i < 5; i++ {
		sem.Run(task)
	}

	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(&running) > int32(concurrency) {
		t.Errorf("expected at most %d running tasks, got %d", concurrency, running)
	}

	sem.Wait()
	if atomic.LoadInt32(&running) != 0 {
		t.Errorf("expected 0 running tasks, got %d", running)
	}
}

func TestSemaphoreWait(t *testing.T) {
	sem := NewSemaphore(1)
	var completed int32

	task := func() {
		time.Sleep(50 * time.Millisecond)
		atomic.AddInt32(&completed, 1)
	}

	for i := 0; i < 3; i++ {
		sem.Run(task)
	}

	sem.Wait()
	if atomic.LoadInt32(&completed) != 3 {
		t.Errorf("expected 3 completed tasks, got %d", completed)
	}
}
