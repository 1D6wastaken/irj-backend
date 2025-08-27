package utils

import (
	"context"
	"errors"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var ErrTimeoutWait = errors.New("timeout has been reached before wait is completed")

type AppStopper struct {
	noCopy noCopy

	cancel context.CancelFunc
	done   <-chan struct{}
	wg     *sync.WaitGroup
}

func NewAppStopper(ctx context.Context) *AppStopper {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	return &AppStopper{
		noCopy: noCopy{},
		cancel: cancel,
		done:   ctx.Done(),
		wg:     &sync.WaitGroup{},
	}
}

func (s *AppStopper) Hold(delta int) {
	s.wg.Add(delta)
}

func (s *AppStopper) Release() {
	s.wg.Done()
}

func (s *AppStopper) Wait(timeout time.Duration) error {
	termination := make(chan struct{})

	go func() {
		defer close(termination)

		s.wg.Wait()
	}()

	timer := time.NewTimer(timeout)

	select {
	case <-termination:
		timer.Stop()

		return nil
	case <-timer.C:
		return ErrTimeoutWait
	}
}

func (s *AppStopper) Cancel() {
	s.cancel()
}

func (s *AppStopper) Done() <-chan struct{} {
	return s.done
}
