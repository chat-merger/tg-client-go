package app

import (
	"context"
	"sync"
)

type application struct {
	errCh      chan<- error
	wg         *sync.WaitGroup // for indicate all things (controllers, servers, handlers...) will stopped
	cancelFunc context.CancelFunc
	ctx        context.Context
}

func newApplication(ctx context.Context) (*application, <-chan error) {
	errCh := make(chan error)
	ctx, cancelFunc := context.WithCancel(ctx)
	return &application{
		errCh:      errCh,
		wg:         new(sync.WaitGroup),
		cancelFunc: cancelFunc,
		ctx:        ctx,
	}, errCh
}
