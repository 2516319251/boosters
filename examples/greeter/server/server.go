package main

import (
	"github.com/2516319251/boosters/examples/greeter/conf"

	grpc2 "google.golang.org/grpc"

	"github.com/2516319251/boosters/internal/testdata/greeter"
	"github.com/2516319251/boosters/transport/grpc"
)

func NewServer(c *conf.Grpc) *grpc.Server {
	gs := grpc2.NewServer()
	greeter.RegisterGreeterServer(gs, &Service{})

	opts := []grpc.ServerOption{
		grpc.Network("tcp"),
		grpc.Address(c.Addr),
	}
	srv := grpc.NewServer(gs, opts...)

	return srv
}
