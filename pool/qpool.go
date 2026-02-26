package pool

import (
	"context"
	"sync"
)

type ProcessRoutine[T any] func(context.Context, int, T)

type QPool[T any] struct {
	tasks  chan T
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	once   sync.Once
	pushWg sync.WaitGroup
}

func NewQPool[T any](workers int, routine ProcessRoutine[T]) *QPool[T] {
	ip := &QPool[T]{
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

func (ip *QPool[T]) Push(task T) error {
	ip.pushWg.Add(1)
	defer ip.pushWg.Done()

	select {
	case <-ip.ctx.Done():
		return ip.ctx.Err()
	case ip.tasks <- task:
		return nil
	}
}

func (ip *QPool[T]) Stop() {
	ip.once.Do(func() {
		ip.cancel()
		ip.pushWg.Wait()
		close(ip.tasks)
		ip.wg.Wait()
	})
}
