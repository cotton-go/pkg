package limiter

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	limiter := NewLimiter(3) // 限制为 3

	for i := 0; i <= 10; i++ {
		go limiter.ExecuteWithTicket(func(t int) {
			st := time.Now()
			defer func() {
				fmt.Println("do something end", "t", t, time.Since(st))
			}()

			// 生成一个随机数
			rand.Seed(time.Now().UnixNano())
			n := 2 + rand.Intn(3)

			// do something
			fmt.Println("do something", "t", t)
			time.Sleep(time.Duration(n) * time.Second)
		})
	}

	time.Sleep(30 * time.Second)
}
