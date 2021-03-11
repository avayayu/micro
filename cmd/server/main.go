package main

import (
	"go.uber.org/zap"
	"gogs.bfr.com/zouhy/micro/micro"
	"gogs.bfr.com/zouhy/micro/net/http"
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
