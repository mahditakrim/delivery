package main

import (
	"context"
	"sync"
)

type workerPool struct {
	wpChan chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func newWorkerPool(size int) *workerPool {

	ctx, cancel := context.WithCancel(context.Background())

	return &workerPool{
		wpChan: make(chan struct{}, size),
		ctx:    ctx,
		cancel: cancel,
		wg:     sync.WaitGroup{},
	}
}

func (wp *workerPool) launch(f func()) {

	select {
	case <-wp.ctx.Done():
		return
	case wp.wpChan <- struct{}{}:
		wp.wg.Add(1)
		go func() {
			defer func() {
				<-wp.wpChan
				wp.wg.Done()
			}()
			f()
		}()
	}
}

func (wp *workerPool) shutdown() {

	wp.cancel()
	wp.wg.Wait()
}
