package pkg

import (
	"sync"
)

func ParallelRun(funcs ...func() (any, error)) ([]any, error) {
	results := make([]any, len(funcs))
	var firstErr error
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(len(funcs))
	for i, fn := range funcs {
		i, fn := i, fn // capture
		go func() {
			defer wg.Done()
			res, err := fn()
			mu.Lock()
			defer mu.Unlock()
			if err != nil && firstErr == nil {
				firstErr = err
			}
			results[i] = res
		}()
	}
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}
	return results, nil
}
