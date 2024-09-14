package pool

import (
	"sync"
)

type Pool interface {
	Reset()
}

type pool[T any] struct {
	p   *sync.Pool
	opt *Option
}

func NewPool[T any](opts ...Options) *pool[T] {
	p := &pool[T]{
		opt: newOptions(),
	}
	p.p = &sync.Pool{
		New: func() any {
			// atomic.AddInt32(&p.opt.count, 1)
			return new(T)
		},
	}

	return p
}

func (p *pool[T]) Get() *T {
	// fmt.Println("get count", atomic.LoadInt32(&p.opt.count))
	value, _ := p.p.Get().(*T)
	return value
}

func (p *pool[T]) Put(v *T) {
	// fmt.Println("put count", atomic.LoadInt32(&p.opt.count))
	p.p.Put(v)
}
