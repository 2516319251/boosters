package data

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"

	"github.com/2516319251/boosters/registry"
	"github.com/2516319251/boosters/registry/etcd"
)

func NewClient() *clientv3.Client {
	cfg := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}

	return client
}

func NewRegistry(client *clientv3.Client) registry.Registrar {
	opts := []etcd.Option{
		etcd.RegisterTTL(time.Hour),
		etcd.Namespace("discovery://"),
	}
	return etcd.New(client, opts...)
}

func NewDiscovery(client *clientv3.Client) registry.Discovery {
	opts := []etcd.Option{
		etcd.RegisterTTL(time.Hour),
		etcd.Namespace("discovery://"),
	}
	return etcd.New(client, opts...)
}
