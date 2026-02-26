package pool

import (
	"context"
	"sync"
)

type ProcessRoutine[T any] func(context.Context, int, T)

type Pool[T any] struct {
	tasks  chan T
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	once   sync.Once
}

func NewPool[T any](workers int, routine ProcessRoutine[T]) *Pool[T] {
	ip := &Pool[T]{
		tasks: make(chan T, workers),
	}

	ip.wg.Add(workers)
	ip.ctx, ip.cancel = context.WithCancel(context.Background())
	for i := 0; i < workers; i++ {
		i := i
		go func() {
			defer ip.wg.Done()

			for task := range ip.tasks {
				routine(ip.ctx, i, task)
			}
		}()
	}

	return ip
}

func (ip *Pool[T]) PushContext(ctx context.Context, task T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case ip.tasks <- task:
		return nil
	}
}

func (ip *Pool[T]) Stop() {
	ip.once.Do(func() {
		ip.cancel()
		close(ip.tasks)
		ip.wg.Wait()
	})
}
