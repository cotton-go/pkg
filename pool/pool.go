package pool

import "sync"

type Pool interface {
	Close() error
}

type pool[T Pool] struct {
	p   *sync.Pool
	opt *Option
}

func NewPool[T Pool](opts ...Options) *pool[T] {
	p := &pool[T]{
		opt: newOptions(),
		p: &sync.Pool{
			New: func() any { return new(T) },
		},
	}

	return p
}

func (p *pool[T]) Get() *T {
	value, _ := p.p.Get().(*T)
	return value
}

func (p *pool[T]) Put(v *T) {
	p.p.Put(v)
}
