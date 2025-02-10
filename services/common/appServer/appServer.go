package appserver

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
)

type server interface {
	Close() error
	Serve(net.Listener) error
}

type AppServer struct {
	started atomic.Bool

	ctx context.Context

	server server

	port int

	errChan chan error

	closers []io.Closer
}

// addr is a network address that must match the form "host:port"
func NewAppServer(ctx context.Context, server server, port int) *AppServer {
	return &AppServer{
		ctx:     ctx,
		port:    port,
		server:  server,
		errChan: make(chan error),
	}
}

func (a *AppServer) WithClosers(closers []io.Closer) *AppServer {
	a.closers = closers
	return a
}

// start listen
func (a *AppServer) Start() error {
	if a.started.Swap(true) {
		return ErrAlreadyStarted
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return err
	}

	errChan := make(chan error)

	go func() {
		defer close(errChan)

		err := a.server.Serve(l)
		if err != nil {
			errChan <- err
		}
	}()

	go func(ctx context.Context) {
		defer func() {
			err := l.Close()
			if err != nil {
				a.errChan <- err
			}

			err = a.server.Close()
			if err != nil {
				a.errChan <- err
			}

			close(a.errChan)
		}()

		select {
		case <-ctx.Done():
			return
		case err := <-errChan:
			a.errChan <- err
			return
		}
	}(a.ctx)

	return nil
}

// waiting when all goroutines is done and return serve errors
func (a *AppServer) wait() []error {
	errs := make([]error, 0)

	for err := range a.errChan {
		errs = append(errs, err)
	}

	return errs
}

// waiting when all goroutines is done and return close and serve erros
func (a *AppServer) WaitAndClose() []error {
	errs := a.wait()

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	wg.Add(len(a.closers))

	for _, closer := range a.closers {
		go func(closer io.Closer) {
			defer wg.Done()

			err := closer.Close()
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}(closer)
	}

	wg.Wait()

	return errs
}
