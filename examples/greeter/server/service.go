package main

import (
	"context"
	"fmt"
	"log"

	"github.com/2516319251/boosters/internal/testdata/greeter"
)

type Service struct {
	greeter.UnimplementedGreeterServer
}

func (svc *Service) SayHello(ctx context.Context, r *greeter.HelloRequest) (*greeter.HelloReply, error) {
	log.Printf("[Received] name: %v", r.GetName())
	return &greeter.HelloReply{Message: fmt.Sprintf("Hello, %s!", r.GetName())}, nil
}
