package conc

import (
	"fmt"
	"sync"
)

// Semaphore is a helper to run multiple goroutines with limited concurrency
type Semaphore struct {
	wg sync.WaitGroup
	ch chan int
}

func NewSemaphore(concurrency int) *Semaphore {
	return &Semaphore{
		ch: make(chan int, concurrency),
	}
}

func (r *Semaphore) Run(task func()) {
	r.wg.Add(1)
	r.ch <- 1
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
			<-r.ch
			r.wg.Done()
		}()

		task()
	}()
}

// Wait waits for all tasks to finish
func (r *Semaphore) Wait() {
	r.wg.Wait()
}
