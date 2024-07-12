package main

import (
	"errors"
	"fmt"

	"github.com/cotton-go/pkg/slicex"
)

func main() {
	for i := range slicex.EmptyArray(10) {
		fmt.Println("i =", i)
	}

	hook := func(i int) (bool, error) {
		if i == 2 {
			// 中途执行退出
			return true, nil
		}

		if i == 3 {
			// 遇到问题中途退出
			return false, errors.New("遇到错误了")
		}

		return false, nil
	}

	if err := slicex.For(5, hook); err != nil {
		// 遇到问题中途退出, 如果中途没有遇到问题也退出了循环则不会进入这个 if
		panic(err)
	}
}
