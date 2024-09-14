package pool

type Option struct {
	Max int
}

type Options func(*Option)

func newOptions(opts ...Options) *Option {
	opt := &Option{}
	for _, o := range opts {
		o(opt)
	}
	return opt
}
