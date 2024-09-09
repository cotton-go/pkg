package main

import (
	"context"
	"fmt"

	"github.com/cotton-go/pkg/signal"
)

func main() {
	fmt.Println("Hello, World!")

	ctx := context.Background()
	signal.Shutdown(ctx, func() {
		fmt.Println("Shutdown")
	})
}
