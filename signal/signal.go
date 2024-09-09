package signal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Hook func()

func Signals() []os.Signal {
	return []os.Signal{
		syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGUSR1, syscall.SIGUSR2,
		os.Kill, os.Interrupt,
	}
}

func Shutdown(ctx context.Context, hooks ...Hook) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, Signals()...)
	fn := func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("panic: %v\n", e)
			}
		}()

		for _, hook := range hooks {
			hook()
		}
	}

	select {
	case <-ctx.Done():
		fn()
		os.Exit(0)
	case <-done:
		fn()
		os.Exit(0)
	}
}
