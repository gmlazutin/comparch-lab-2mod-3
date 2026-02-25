package pool

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrPoolIsClosing = errors.New("pool is in closing state")
)

type ProcessRoutine[T any] func(context.Context, int, T)

type Pool[T any] struct {
	tasks   chan T
	ctx     context.Context
	cancel  context.CancelFunc
	closing atomic.Bool
	wg      sync.WaitGroup
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

			for {
				select {
				case task := <-ip.tasks:
					routine(ip.ctx, i, task)
				case <-ip.ctx.Done():
					return
				}
			}
		}()
	}

	return ip
}

func (ip *Pool[T]) Push(task T) error {
	if !ip.closing.Load() {
		ip.tasks <- task
		return nil
	}

	return ErrPoolIsClosing
}

func (ip *Pool[T]) Stop() {
	if !ip.closing.CompareAndSwap(false, true) {
		return
	}
	close(ip.tasks)
	ip.cancel()
	ip.wg.Wait()
}
