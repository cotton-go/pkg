package pool

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

type item struct {
	Name string
}

func (i *item) Reset() {
	i.Name = ""
}

func TestPool(t *testing.T) {
	var wg sync.WaitGroup
	p := NewPool[item]()
	for i := range 1000 {
		if i%20 == 0 {
			time.Sleep(time.Second)
		}

		wg.Add(1)
		pi := p.Get()
		pi.Name = "test:" + strconv.Itoa(i)
		go func(pi *item) {
			defer wg.Done()
			defer p.Put(pi)
			time.Sleep(time.Second)
			t.Log("pi", pi, "pi.name", pi.Name)
		}(pi)
	}

	wg.Wait()
}
