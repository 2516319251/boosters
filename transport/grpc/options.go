package grpc

type ServerOption func(o *options)

type options struct {
	network string
	address string
}

func Network(network string) ServerOption {
	return func(o *options) {
		o.network = network
	}
}

func Address(address string) ServerOption {
	return func(o *options) {
		o.address = address
	}
}
