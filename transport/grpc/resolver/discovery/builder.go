package discovery

import (
	"context"
	"errors"
	"strings"
	"time"

	"google.golang.org/grpc/resolver"

	"github.com/2516319251/boosters/registry"
)

const name = "discovery"

type builder struct {
	opts       option
	discoverer registry.Discovery
}

// NewBuilder creates a builder which is used to factory registry resolvers.
func NewBuilder(d registry.Discovery, opts ...Option) resolver.Builder {
	o := option{
		timeout:          time.Second * 10,
		insecure:         false,
		debugLogDisabled: false,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &builder{
		opts:       o,
		discoverer: d,
	}
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	res := &struct {
		watcher registry.Watcher
		err     error
	}{}

	done := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		w, err := b.discoverer.Watch(ctx, strings.TrimPrefix(target.URL.Path, "/"))
		res.watcher = w
		res.err = err
		close(done)
	}()

	var err error
	select {
	case <-done:
		err = res.err
	case <-time.After(b.opts.timeout):
		err = errors.New("discovery create watcher overtime")
	}

	if err != nil {
		cancel()
		return nil, err
	}

	r := &discoveryResolver{
		watcher:          res.watcher,
		conn:             cc,
		ctx:              ctx,
		cancel:           cancel,
		insecure:         b.opts.insecure,
		debugLogDisabled: b.opts.debugLogDisabled,
	}
	go r.watch()

	return r, nil
}

// Scheme return scheme of discovery
func (*builder) Scheme() string {
	return name
}
