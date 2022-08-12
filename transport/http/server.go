package http

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/2516319251/boosters/transport"
)

type Server struct {
	server *http.Server
}

func NewServer(server *http.Server) transport.Server {
	return &Server{
		server: server,
	}
}

func (srv *Server) Start(_ context.Context) error {
	log.Printf("[HTTP] server listening on: %s\n", srv.server.Addr)

	// 启动 http 服务
	err := srv.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (srv *Server) Stop(ctx context.Context) error {
	log.Println("[HTTP] server stop")

	// 停止 http 服务
	return srv.server.Shutdown(ctx)
}
