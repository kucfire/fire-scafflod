package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

type Hook interface {
	WithSignals(signals ...syscall.Signal) Hook

	Close(funcs ...func())
}

type hook struct {
	ctx chan os.Signal
}

func New() Hook {
	newHook := &hook{
		ctx: make(chan os.Signal, 1),
	}

	return newHook
}

func (h *hook) WithSignals(signals ...syscall.Signal) Hook {
	for _, s := range signals {
		signal.Notify(h.ctx, s)
	}

	return h
}

func (h *hook) Close(funcs ...func()) {
	select {
	case <-h.ctx:
	}

	signal.Stop(h.ctx)

	for _, f := range funcs {
		f()
	}
}
