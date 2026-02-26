package pool

import (
	"context"
	"sync"
)

type ProcessRoutine[T any] func(context.Context, int, T)

type QPool[T any] struct {
	tasks chan T
	wg    sync.WaitGroup
	ctx   context.Context
}

func NewQPool[T any](ctx context.Context, workers int, routine ProcessRoutine[T]) *QPool[T] {
	ip := &QPool[T]{
		tasks: make(chan T, workers),
		ctx:   ctx,
	}

	ip.wg.Add(workers)
	for i := 0; i < workers; i++ {
		i := i
		go func() {
			defer ip.wg.Done()

			for task := range ip.tasks {
				if ctx.Err() != nil {
					break
				}
				routine(ctx, i, task)
			}
		}()
	}

	return ip
}

func (ip *QPool[T]) PushContext(ctx context.Context, task T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case ip.tasks <- task:
		return nil
	}
}

func (ip *QPool[T]) Push(task T) error {
	return ip.PushContext(context.Background(), task)
}

func (ip *QPool[T]) WaitDone() error {
	close(ip.tasks)
	ip.wg.Wait()

	return ip.ctx.Err()
}
