package rpc

import (
	"context"
	"math"
	"net"
	"sync"
	"time"

	nmd "github.com/avayayu/micro/net/metadata"
	"github.com/avayayu/micro/net/trace"
	ztime "github.com/avayayu/micro/time"
	"github.com/pkg/errors"
	"github.com/siddontang/go/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
)

var _abortIndex int8 = math.MaxInt8 / 2

type RpcServerConfig struct {
	Addr              string
	Timeout           ztime.Duration
	IdleTimeout       ztime.Duration
	MaxLifeTime       ztime.Duration
	ForceCloseWait    ztime.Duration
	KeepAliveInterval ztime.Duration
	KeepAliveTimeout  ztime.Duration
}

func DefaultRpcServerConfig() *RpcServerConfig {
	return &RpcServerConfig{
		Addr:              "0.0.0.0:9000",
		Timeout:           ztime.Duration(time.Second),
		IdleTimeout:       ztime.Duration(time.Second * 180),
		MaxLifeTime:       ztime.Duration(time.Hour * 2),
		ForceCloseWait:    ztime.Duration(time.Second * 20),
		KeepAliveInterval: ztime.Duration(time.Second * 60),
		KeepAliveTimeout:  ztime.Duration(time.Second * 20),
	}
}

type RpcServer struct {
	config   *RpcServerConfig
	mutex    sync.RWMutex
	server   *grpc.Server
	handlers []grpc.UnaryServerInterceptor
	logger   *zap.Logger
}

// handle return a new unary server interceptor for OpenTracing\Logging\LinkTimeout.
func (s *RpcServer) handle() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var (
			cancel func()
			addr   string
		)
		s.mutex.RLock()
		config := s.config
		s.mutex.RUnlock()
		// get derived timeout from grpc context,
		// compare with the warden configured,
		// and use the minimum one
		timeout := time.Duration(config.Timeout)
		if dl, ok := ctx.Deadline(); ok {
			ctimeout := time.Until(dl)
			if ctimeout-time.Millisecond*20 > 0 {
				ctimeout = ctimeout - time.Millisecond*20
			}
			if timeout > ctimeout {
				timeout = ctimeout
			}
		}
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()

		// get grpc metadata(trace & remote_ip & color)
		var t trace.Trace
		cmd := nmd.MD{}
		if gmd, ok := metadata.FromIncomingContext(ctx); ok {
			t, _ = trace.Extract(trace.GRPCFormat, gmd)
			for key, vals := range gmd {
				if nmd.IsIncomingKey(key) {
					cmd[key] = vals[0]
				}
			}
		}
		if t == nil {
			t = trace.New(args.FullMethod)
		} else {
			t.SetTitle(args.FullMethod)
		}

		if pr, ok := peer.FromContext(ctx); ok {
			addr = pr.Addr.String()
			t.SetTag(trace.String(trace.TagAddress, addr))
		}
		defer t.Finish(&err)

		// use common meta data context instead of grpc context
		ctx = nmd.NewContext(ctx, cmd)
		ctx = trace.NewContext(ctx, t)

		resp, err = handler(ctx, req)
		return resp, err
	}
}

// SetConfig hot reloads server config
func (s *RpcServer) SetConfig(config *RpcServerConfig) (err error) {
	if config == nil {
		panic("config not not be nill")
	}
	if config.Timeout <= 0 {
		config.Timeout = ztime.Duration(time.Second)
	}
	if config.IdleTimeout <= 0 {
		config.IdleTimeout = ztime.Duration(time.Second * 60)
	}
	if config.MaxLifeTime <= 0 {
		config.MaxLifeTime = ztime.Duration(time.Hour * 2)
	}
	if config.ForceCloseWait <= 0 {
		config.ForceCloseWait = ztime.Duration(time.Second * 20)
	}
	if config.KeepAliveInterval <= 0 {
		config.KeepAliveInterval = ztime.Duration(time.Second * 60)
	}
	if config.KeepAliveTimeout <= 0 {
		config.KeepAliveTimeout = ztime.Duration(time.Second * 20)
	}
	if config.Addr == "" {
		config.Addr = "0.0.0.0:9000"
	}

	s.mutex.Lock()
	s.config = config
	s.mutex.Unlock()
	return nil
}

// NewServer returns a new blank Server instance with a default server interceptor.
func NewServer(config *RpcServerConfig, opt ...grpc.ServerOption) (s *RpcServer) {
	if config == nil {
		config = DefaultRpcServerConfig()
	}
	s = new(RpcServer)
	if err := s.SetConfig(config); err != nil {
		panic(errors.Errorf("warden: set config failed!err: %s", err.Error()))
	}
	keepParam := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(s.config.IdleTimeout),
		MaxConnectionAgeGrace: time.Duration(s.config.ForceCloseWait),
		Time:                  time.Duration(s.config.KeepAliveInterval),
		Timeout:               time.Duration(s.config.KeepAliveTimeout),
		MaxConnectionAge:      time.Duration(s.config.MaxLifeTime),
	})
	opt = append(opt, keepParam, grpc.UnaryInterceptor(s.interceptor))
	s.server = grpc.NewServer(opt...)
	s.Use(s.recovery(), s.handle(), s.stats())
	// s.Use(ratelimiter.New(nil).Limit())
	return
}

// interceptor is a single interceptor out of a chain of many interceptors.
// Execution is done in left-to-right order, including passing of context.
// For example ChainUnaryServer(one, two, three) will execute one before two before three, and three
// will see context changes of one and two.
func (s *RpcServer) interceptor(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var (
		i     int
		chain grpc.UnaryHandler
	)

	n := len(s.handlers)
	if n == 0 {
		return handler(ctx, req)
	}

	chain = func(ic context.Context, ir interface{}) (interface{}, error) {
		if i == n-1 {
			return handler(ic, ir)
		}
		i++
		return s.handlers[i](ic, ir, args, chain)
	}

	return s.handlers[0](ctx, req, args, chain)
}

// Use attachs a global inteceptor to the server.
// For example, this is the right place for a rate limiter or error management inteceptor.
func (s *RpcServer) Use(handlers ...grpc.UnaryServerInterceptor) *RpcServer {
	finalSize := len(s.handlers) + len(handlers)
	if finalSize >= int(_abortIndex) {
		panic("warden: server use too many handlers")
	}
	mergedHandlers := make([]grpc.UnaryServerInterceptor, finalSize)
	copy(mergedHandlers, s.handlers)
	copy(mergedHandlers[len(s.handlers):], handlers)
	s.handlers = mergedHandlers
	return s
}

// Server return the grpc server for registering service.
func (s *RpcServer) Server() *grpc.Server {
	return s.server
}

// Run create a tcp listener and start goroutine for serving each incoming request.
// Run will return a non-nil error unless Stop or GracefulStop is called.
func (s *RpcServer) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		err = errors.WithStack(err)
		log.Error("failed to listen: %v", err)
		return err
	}
	reflection.Register(s.server)
	return s.Serve(lis)
}

// Start create a new goroutine run server with configured listen addr
// will panic if any error happend
// return server itself
func (s *RpcServer) Start() (*RpcServer, error) {
	_, err := s.startWithAddr()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *RpcServer) startWithAddr() (net.Addr, error) {
	lis, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		return nil, err
	}
	log.Info("warden: start grpc listen addr: %v", lis.Addr())
	reflection.Register(s.server)
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return lis.Addr(), nil
}

// Serve accepts incoming connections on the listener lis, creating a new
// ServerTransport and service goroutine for each.
// Serve will return a non-nil error unless Stop or GracefulStop is called.
func (s *RpcServer) Serve(lis net.Listener) error {
	return s.server.Serve(lis)
}


func (s *RpcServer) Shutdown(ctx context.Context) (err error) {
	ch := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(ch)
	}()
	select {
	case <-ctx.Done():
		s.server.Stop()
		err = ctx.Err()
	case <-ch:
	}
	return
}

func (s *RpcServer) RegisterHandler()
