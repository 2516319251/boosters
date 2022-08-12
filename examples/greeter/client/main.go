package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/2516319251/boosters/examples/greeter/data"
	"github.com/2516319251/boosters/internal/testdata/greeter"
	"github.com/2516319251/boosters/transport/grpc/resolver/discovery"
)

func main() {
	client := data.NewClient()
	defer client.Close()

	// Set up a connection to the server.
	opts := []grpc.DialOption{
		grpc.WithResolvers(discovery.NewBuilder(data.NewDiscovery(client))),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial("discovery:///greeter", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// With timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Contact the server and print out its response.
	c := greeter.NewGreeterClient(conn)
	r, err := c.SayHello(ctx, &greeter.HelloRequest{Name: "ZEE"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

	select {}
}
