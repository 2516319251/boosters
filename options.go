package boosters

import (
	"context"
	"net/url"
	"os"

	"github.com/2516319251/boosters/registry"
	"github.com/2516319251/boosters/transport"
)

type Option func(o *options)

type options struct {
	id, name, version string
	endpoints         []*url.URL

	ctx  context.Context
	sigs []os.Signal

	registrar registry.Registrar
	servers   []transport.Server
}

func ID(id string) Option {
	return func(o *options) { o.id = id }
}

func Name(name string) Option {
	return func(o *options) { o.name = name }
}

func Version(version string) Option {
	return func(o *options) { o.version = version }
}

func Endpoint(endpoints ...*url.URL) Option {
	return func(o *options) { o.endpoints = endpoints }
}

func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

func Signal(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

func Registrar(r registry.Registrar) Option {
	return func(o *options) { o.registrar = r }
}

func Server(srv ...transport.Server) Option {
	return func(o *options) { o.servers = srv }
}
