package conc

import (
	"fmt"
	"sync"
)

// Each iterates over the given items and executes the given handler in parallel
func Each[Input any](num int, items []Input, handler func(Input) error) error {
	wg := sync.WaitGroup{}
	ch := make(chan struct{}, num)
	mu := sync.Mutex{}
	var errs []error

	wg.Add(len(items))

	for i := range items {
		ch <- struct{}{}
		item := items[i]
		go func() {
			defer func() {
				if r := recover(); r != nil {
					mu.Lock()
					errs = append(errs, fmt.Errorf("panic %v", r))
					mu.Unlock()
				}
				<-ch
				wg.Done()
			}()

			err := handler(item)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if len(errs) > 0 {
		return MultiError(errs)
	}

	return nil
}

type MultiError []error

func (me MultiError) Error() string {
	return fmt.Sprintf("%v", []error(me))
}
