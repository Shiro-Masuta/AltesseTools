package images

import (
	"fmt"
	"runtime"
	"sync"
)

func ParallelConvert[T any](items []T, fn func(int, T) ([]byte, error)) ([][]byte, error) {
	// Limite de workers = nb CPU (max 8 pour ne pas saturer la machine)
	workers := runtime.NumCPU()
	if workers > 8 {
		workers = 8
	}

	results := make([][]byte, len(items))
	var firstError error
	var mu sync.Mutex

	var wg sync.WaitGroup
	tasks := make(chan int)

	// Création des workers
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range tasks {
				res, err := fn(i, items[i])

				mu.Lock()
				if err != nil && firstError == nil {
					firstError = fmt.Errorf("index %d : %w", i, err)
				}
				results[i] = res
				mu.Unlock()
			}
		}()
	}

	// Envoi des indices dans le channel
	go func() {
		for i := range items {
			tasks <- i
		}
		close(tasks)
	}()

	wg.Wait()

	if firstError != nil {
		return nil, firstError
	}
	return results, nil
}
