package main

import (
	"github.com/avayayu/micro/micro"
	"github.com/avayayu/micro/net/http"
	"go.uber.org/zap"
)

const _module = "test"

type Server struct {
	engine *http.Engine
	logger *zap.Logger
}

func (server *Server) ServiceName() string {

	return _module
}

func (server *Server) MicroConfig() *http.ServerConfig {

	return nil
}

func (server *Server) HttpEngine(serverConfig *http.ServerConfig) *http.Engine {
	server.engine = http.NewServer(serverConfig)
	return server.engine
}

func (server *Server) Logger() *zap.Logger {
	return server.logger
}

func (server *Server) Close() {

}

func (server *Server) SetupHandler() {}

func NewServer() *Server {
	server := &Server{}

	return server
}

func main() {
	server := NewServer()

	micro.RunService(server)
}
