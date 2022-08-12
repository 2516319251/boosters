package grpc

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	opts     options
	server   *grpc.Server
	health   *health.Server
	listener net.Listener
	err      error
}

func NewServer(server *grpc.Server, opts ...ServerOption) *Server {
	o := options{
		network: "tcp",
		address: ":0",
	}

	for _, opt := range opts {
		opt(&o)
	}

	srv := &Server{
		opts:   o,
		server: server,
		health: health.NewServer(),
	}

	// 设置 grpc 的监听地址
	srv.err = srv.listen()
	grpc_health_v1.RegisterHealthServer(srv.server, srv.health)

	return srv
}

func (srv *Server) listen() error {
	// 设置监听地址
	if srv.listener == nil {
		lis, err := net.Listen(srv.opts.network, srv.opts.address)
		if err != nil {
			return err
		}
		srv.listener = lis
	}
	return nil
}

func (srv *Server) Start(_ context.Context) error {
	log.Printf("[GRPC] server listening on: %s", srv.listener.Addr().String())

	// 如果存在错误
	if srv.err != nil {
		return srv.err
	}

	// 启动 grpc 服务
	srv.health.Resume()
	return srv.server.Serve(srv.listener)
}

func (srv *Server) Stop(_ context.Context) error {
	log.Printf("[GRPC] server stop")

	// 停止 grpc 服务
	srv.health.Shutdown()
	srv.server.GracefulStop()

	return nil
}
