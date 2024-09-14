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
	return []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL}
}

func shutdownctx(ctx context.Context) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, Signals()...)
	select {
	// wait on kill signal
	case <-ctx.Done():
	// wait on context cancel
	case <-done:
	}
}

func Shutdown(ctx context.Context, hooks ...Hook) {
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

	go shutdownctx(ctx)
	shutdownctx(ctx)
	fn()
}
