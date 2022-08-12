package config

type options struct {
	source  Source
}

type Option func(o *options)

func WithSource(s Source) Option {
	return func(o *options) {
		o.source = s
	}
}
